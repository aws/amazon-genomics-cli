import logging
from abc import ABC, abstractmethod
from typing import Optional

from amazon_genomics.util.method_logger import logged

from rest_api.models import RunLog, RunId, RunStatus, ServiceInfo, RunListResponse

DEFAULT_LOGGER_FORMAT = "[%(asctime)-15s] [%(name)s] [%(levelname)s] %(message)s"


class AbstractWESAdapter(ABC):
    """
    Abstract WES Interface
    """

    def __init__(self, logger=None):
        # Create logger
        if logger:
            self.logger = logger
        else:
            logging.basicConfig(format=DEFAULT_LOGGER_FORMAT)
            self.logger = logging.getLogger(type(self).__name__)
            self.logger.setLevel(logging.DEBUG)

        # Shut off the logging of every http request by werkzeug
        wz_logger = logging.getLogger("werkzeug")
        wz_logger.disabled = True

    @abstractmethod
    def cancel_run(self, run_id: str) -> Optional[RunId]:
        """Cancel a running workflow.

        # noqa: E501

        :param run_id:
        :type run_id: str

        :rtype: RunId
        """
        pass

    @abstractmethod
    def get_run_log(self, run_id) -> Optional[RunLog]:
        """Get detailed info about a workflow run.

        This endpoint provides detailed information about a given workflow run. The returned result has information about the outputs produced by this workflow (if available), a log object which allows the stderr and stdout to be retrieved, a log array so stderr/stdout for individual tasks can be retrieved, and the overall state of the workflow run (e.g. RUNNING, see the State section). # noqa: E501

        :param run_id:
        :type run_id: str

        :rtype: RunLog
        """
        pass

    @abstractmethod
    def get_run_status(self, run_id) -> Optional[RunStatus]:
        """Get quick status info about a workflow run.

        This provides an abbreviated (and likely fast depending on implementation) status of the running workflow, returning a simple result with the  overall state of the workflow run (e.g. RUNNING, see the State section). # noqa: E501

        :param run_id:
        :type run_id: str

        :rtype: RunStatus
        """
        pass

    @abstractmethod
    def get_service_info(self) -> ServiceInfo:
        """Get information about Workflow Execution Service.

        May include information related (but not limited to) the workflow descriptor formats, versions supported, the WES API versions supported, and information about general service availability. # noqa: E501


        :rtype: ServiceInfo
        """
        pass

    @abstractmethod
    def list_runs(self, page_size=None, page_token=None) -> RunListResponse:
        """List the workflow runs.

        This list should be provided in a stable ordering. (The actual ordering is implementation dependent.) When paging through the list, the client should not make assumptions about live updates, but should assume the contents of the list reflect the workflow list at the moment that the first page is requested.  To monitor a specific workflow run, use GetRunStatus or GetRunLog. # noqa: E501

        :param page_size: OPTIONAL The preferred number of workflow runs to return in a page. If not provided, the implementation should use a default page size. The implementation must not return more items than &#x60;page_size&#x60;, but it may return fewer.  Clients should not assume that if fewer than &#x60;page_size&#x60; items are returned that all items have been returned.  The availability of additional pages is indicated by the value of &#x60;next_page_token&#x60; in the response.
        :type page_size: int
        :param page_token: OPTIONAL Token to use to indicate where to start getting results. If unspecified, return the first page of results.
        :type page_token: str

        :rtype: RunListResponse
        """
        pass

    @abstractmethod
    def run_workflow(
        self,
        workflow_params=None,
        workflow_type=None,
        workflow_type_version=None,
        tags=None,
        workflow_engine_parameters=None,
        workflow_url=None,
        workflow_attachment=None,
    ) -> RunId:
        """Run a workflow.
        This endpoint creates a new workflow run and returns a &#x60;RunId&#x60; to monitor its progress.  The &#x60;workflow_attachment&#x60; array may be used to upload files that are required to execute the workflow, including the primary workflow, tools imported by the workflow, other files referenced by the workflow, or files which are part of the input.  The implementation should stage these files to a temporary directory and execute the workflow from there. These parts must have a Content-Disposition header with a \&quot;filename\&quot; provided for each part.  Filenames may include subdirectories, but must not include references to parent directories with &#39;..&#39; -- implementations should guard against maliciously constructed filenames.  The &#x60;workflow_url&#x60; is either an absolute URL to a workflow file that is accessible by the WES endpoint, or a relative URL corresponding to one of the files attached using &#x60;workflow_attachment&#x60;.  The &#x60;workflow_params&#x60; JSON object specifies input parameters, such as input files.  The exact format of the JSON object depends on the conventions of the workflow language being used.  Input files should either be absolute URLs, or relative URLs corresponding to files uploaded using &#x60;workflow_attachment&#x60;.  The WES endpoint must understand and be able to access URLs supplied in the input.  This is implementation specific.  The &#x60;workflow_type&#x60; is the type of workflow language and must be \&quot;CWL\&quot; or \&quot;WDL\&quot; currently (or another alternative  supported by this WES instance).  The &#x60;workflow_type_version&#x60; is the version of the workflow language submitted and must be one supported by this WES instance.  See the &#x60;RunRequest&#x60; documentation for details about other fields. # noqa: E501

        :rtype: RunId
        """
        pass

    @property
    @abstractmethod
    def workflow_type_versions(self):
        """Workflow type versions supported by an engine"""
        pass

    @property
    def supported_wes_versions(self):
        return ["1.0.0"]

    @logged
    def _is_supported_workflow(
        self, workflow_type: str, workflow_type_version: str
    ) -> bool:
        workflow_type = workflow_type.strip().upper()
        workflow_type_version = workflow_type_version.strip()
        return (
            workflow_type in self.workflow_type_versions
            and workflow_type_version
            in self.workflow_type_versions[workflow_type].workflow_type_version
        )
