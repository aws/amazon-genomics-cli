import os
import json

from flask import request

from rest_api.models import (
    RunId,
    RunListResponse,
    RunLog,
    RunStatus,
    ServiceInfo,
)

from rest_api.exception.Exceptions import InvalidRequestError

from amazon_genomics.util.method_logger import logged
from amazon_genomics.wes.adapters import (
    NextflowWESAdapter,
    CromwellWESAdapter,
    MiniWdlWESAdapter,
)


ENGINE_NAME = os.getenv("ENGINE_NAME")
JOB_QUEUE = os.getenv("JOB_QUEUE")
JOB_DEFINITION = os.getenv("JOB_DEFINITION")
ENGINE_LOG_GROUP = os.getenv("ENGINE_LOG_GROUP")
OUTPUT_DIR_S3_URI = os.getenv("OUTPUT_DIR_S3_URI")

if ENGINE_NAME == "nextflow":
    print("Using Nextflow adapter")
    adapter = NextflowWESAdapter(
        job_queue=JOB_QUEUE,
        job_definition=JOB_DEFINITION,
        engine_log_group=ENGINE_LOG_GROUP,
    )

elif ENGINE_NAME == "cromwell":
    print("Using Cromwell adapter")
    adapter = CromwellWESAdapter()
elif ENGINE_NAME == "miniwdl":
    print("Using MiniWDL adapter")
    adapter = MiniWdlWESAdapter(
        job_queue=JOB_QUEUE,
        job_definition=JOB_DEFINITION,
        output_dir_s3_uri=OUTPUT_DIR_S3_URI,
    )
else:
    raise Exception(f"Unknown engine name `{ENGINE_NAME}`")


@logged
def cancel_run(run_id):  # noqa: E501
    """Cancel a running workflow.

     # noqa: E501

    :param run_id:
    :type run_id: str

    :rtype: RunId
    """
    return adapter.cancel_run(run_id)


@logged
def get_run_log(run_id):  # noqa: E501
    """Get detailed info about a workflow run.

    This endpoint provides detailed information about a given workflow run. The returned result has information about the outputs produced by this workflow (if available), a log object which allows the stderr and stdout to be retrieved, a log array so stderr/stdout for individual tasks can be retrieved, and the overall state of the workflow run (e.g. RUNNING, see the State section). # noqa: E501

    :param run_id:
    :type run_id: str

    :rtype: RunLog
    """
    return adapter.get_run_log(run_id)


@logged
def get_run_status(run_id):  # noqa: E501
    """Get quick status info about a workflow run.

    This provides an abbreviated (and likely fast depending on implementation) status of the running workflow, returning a simple result with the  overall state of the workflow run (e.g. RUNNING, see the State section). # noqa: E501

    :param run_id:
    :type run_id: str

    :rtype: RunStatus
    """
    return adapter.get_run_status(run_id)


@logged
def get_service_info():  # noqa: E501
    """Get information about Workflow Execution Service.

    May include information related (but not limited to) the workflow descriptor formats, versions supported, the WES API versions supported, and information about general service availability. # noqa: E501


    :rtype: ServiceInfo
    """
    return adapter.get_service_info()


@logged
def list_runs(page_size=None, page_token=None):  # noqa: E501
    """List the workflow runs.

    This list should be provided in a stable ordering. (The actual ordering is implementation dependent.) When paging through the list, the client should not make assumptions about live updates, but should assume the contents of the list reflect the workflow list at the moment that the first page is requested.  To monitor a specific workflow run, use GetRunStatus or GetRunLog. # noqa: E501

    :param page_size: OPTIONAL The preferred number of workflow runs to return in a page. If not provided, the implementation should use a default page size. The implementation must not return more items than &#x60;page_size&#x60;, but it may return fewer.  Clients should not assume that if fewer than &#x60;page_size&#x60; items are returned that all items have been returned.  The availability of additional pages is indicated by the value of &#x60;next_page_token&#x60; in the response.
    :type page_size: int
    :param page_token: OPTIONAL Token to use to indicate where to start getting results. If unspecified, return the first page of results.
    :type page_token: str

    :rtype: RunListResponse
    """
    return adapter.list_runs(page_size=page_size, page_token=page_token)


@logged
def run_workflow(
    workflow_params=None,
    workflow_type=None,
    workflow_type_version=None,
    tags=None,
    workflow_engine_parameters=None,
    workflow_url=None,
    workflow_attachment=None,
):  # noqa: E501
    """Run a workflow.

    This endpoint creates a new workflow run and returns a &#x60;RunId&#x60; to monitor its progress.  The &#x60;workflow_attachment&#x60; array may be used to upload files that are required to execute the workflow, including the primary workflow, tools imported by the workflow, other files referenced by the workflow, or files which are part of the input.  The implementation should stage these files to a temporary directory and execute the workflow from there. These parts must have a Content-Disposition header with a \&quot;filename\&quot; provided for each part.  Filenames may include subdirectories, but must not include references to parent directories with &#39;..&#39; -- implementations should guard against maliciously constructed filenames.  The &#x60;workflow_url&#x60; is either an absolute URL to a workflow file that is accessible by the WES endpoint, or a relative URL corresponding to one of the files attached using &#x60;workflow_attachment&#x60;.  The &#x60;workflow_params&#x60; JSON object specifies input parameters, such as input files.  The exact format of the JSON object depends on the conventions of the workflow language being used.  Input files should either be absolute URLs, or relative URLs corresponding to files uploaded using &#x60;workflow_attachment&#x60;.  The WES endpoint must understand and be able to access URLs supplied in the input.  This is implementation specific.  The &#x60;workflow_type&#x60; is the type of workflow language and must be \&quot;CWL\&quot; or \&quot;WDL\&quot; currently (or another alternative  supported by this WES instance).  The &#x60;workflow_type_version&#x60; is the version of the workflow language submitted and must be one supported by this WES instance.  See the &#x60;RunRequest&#x60; documentation for details about other fields. # noqa: E501

    :rtype: RunId
    """

    # parameters are not getting passed through as expected - all string values are set to None.
    # parse the original request object instead
    args = {
        "workflow_params": None,
        "workflow_type": None,
        "workflow_type_version": None,
        "tags": None,
        "workflow_engine_parameters": None,
        "workflow_url": None,
        "workflow_attachment": None,
    }

    for arg in args:

        if arg in ("workflow_attachment"):
            # file lists
            args[arg] = request.files.getlist(arg)
        else:
            args[arg] = request.form.get(arg)

        # set empty values ('', [], etc ...) to None
        if not args[arg]:
            args[arg] = None

        if args[arg] and arg in (
            "tags",
            "workflow_engine_parameters",
            "workflow_params",
        ):
            # json parameters
            try:
                args[arg] = json.loads(args[arg])
            except json.decoder.JSONDecodeError as e:
                raise InvalidRequestError(f"Error processing '{arg}': {e}")

    if not adapter._is_supported_workflow(
        args["workflow_type"], args["workflow_type_version"]
    ):
        msg = f"Unsupported workflow type or version: ({args['workflow_type']}, {args['workflow_type_version']})"
        raise RuntimeError(msg)

    return adapter.run_workflow(**args)
