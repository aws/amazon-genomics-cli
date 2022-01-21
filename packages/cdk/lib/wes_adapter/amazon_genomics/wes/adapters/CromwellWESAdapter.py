import datetime
import json
import os
from os import path
import tempfile
from urllib.parse import urlparse
import zipfile
from typing import Optional
import botocore
import boto3
import requests

from amazon_genomics.wes.adapters import AbstractWESAdapter
from rest_api.exception.Exceptions import InvalidRequestError, InternalServerError
from rest_api.models import (
    ServiceInfo,
    WorkflowTypeVersion,
    RunId,
    RunLog,
    Log,
    RunRequest,
    RunStatus,
    RunListResponse,
    State,
)

appName = "agc"
projectName = os.getenv("PROJECT_NAME", "")
contextName = os.getenv("CONTEXT_NAME", "")
userId = os.getenv("USER_ID", "")
engineServiceName = os.getenv("ENGINE_NAME", "cromwell")
# defaults to older CLI (<= 1.0.1) behavior that used Cloud Map to discover cromwell
engineEndpoint = os.getenv(
    "ENGINE_ENDPOINT",
    f"{engineServiceName}.{projectName}-{contextName}-{userId}.{appName}.amazon.com:8000",
)

WORKFLOW_TYPE_VERSIONS = {"WDL": WorkflowTypeVersion(["1.0", "draft-2"])}

CROMWELL_SERVER = f"http://{engineEndpoint}"
CROMWELL_API_PATH = "api/workflows/v1"
CROMWELL_ENGINE_PATH = "engine/v1"

GET = "GET"
POST = "POST"


class CromwellWESAdapter(AbstractWESAdapter):  # inherit from ABC to enforce interface
    """
    Remote WES adapter that handles WES requests for a WES workflow engine that's
    already running in AWS ECS.
    """

    def __init__(self, logger=None, url_prefix=None, workflow_params=None):
        super().__init__(logger)

        if not url_prefix:
            url_prefix = CROMWELL_SERVER

        if not workflow_params:
            workflow_params = {}

        self.logger.info(f"Initializing: remote server: {url_prefix}")

        self.url_prefix = f"{url_prefix}/{CROMWELL_API_PATH}"
        self.health_check_url = f"{url_prefix}/{CROMWELL_ENGINE_PATH}/status"
        self.workflow_params = workflow_params

    def cancel_run(self, run_id) -> Optional[RunId]:
        """Cancel a running workflow in the remote WES
        workflow engine by calling WES abort REST api
        :param run_id:
        :type run_id: str

        :rtype: run_id: str
        """
        url = self._server_path(run_id, "abort")

        response = requests.request(POST, url)

        return RunId(response.json()["id"])

    def get_run_log(self, run_id) -> Optional[RunLog]:
        """Get detailed info about a workflow run.
        Information is retrieved from workflow metadata and output
        from remote WES workflow engine.
        :param run_id:
        :type run_id: str

        :rtype: run_log_dict: dict
        """
        # Endpoint to get metadata for a specified workflow
        metadata_url = self._server_path(run_id, "metadata")

        # Endpoint to get output for a specified workflow
        outputs_url = self._server_path(run_id, "outputs")

        metadata = requests.request(GET, metadata_url).json()
        outputs = requests.request(GET, outputs_url).json()
        self.logger.info(f"get_run_log metadata: {metadata}")
        self.logger.info(f"get_run_log outputs: {outputs}")
        run_log_dict = self._build_run_log_dict_(metadata, outputs)
        return self._build_run_log_model_(run_log_dict)

    def get_run_status(self, run_id) -> Optional[RunStatus]:
        """Get quick status info about a workflow run.
        Status retrieves current state via the remote WES server's REST API
        :param run_id:
        :type run_id: str

        :rtype: run_status_dict: dict
        """

        url = self._server_path(run_id, "status")

        self.logger.info(f"make GET request to remote WES server at {url}")
        response = requests.request(GET, url)
        self.logger.info(
            f"GET request to remote WES server returns {response} ({response.text})"
        )

        return RunStatus(
            run_id=run_id, state=self._translate_from_response_to_state_(response)
        )

    def get_service_info(self):
        """Get information related (but not limited to)
        the workflow descriptor formats, versions supported,
        the WES API versions supported,
        and information about general service availability.

        :rtype: serviceInfo_dict: dict
        """
        self.logger.info(f"GET_SERVICE_INFO")

        # Check if the WES service engine is healthy
        is_healthy = True
        try:
            self._check_if_wes_service_healthy_()
        except:
            is_healthy = False

        get_service_info_response = ServiceInfo(
            workflow_type_versions=self.workflow_type_versions,
            supported_wes_versions=self.supported_wes_versions,
            tags={
                "name": "remote_cromwell_wes_adapter",
                "description": "WES adapter for Cromwell workflow engine service.",
                "updated_at": datetime.datetime.now(),
                "cromwell_service_health": str(is_healthy),
            },
        )
        return get_service_info_response

    # TODO : Support pagination
    def list_runs(self, page_size=None, page_token=None) -> RunListResponse:
        """List the workflow runs.

        :param page_size: OPTIONAL The preferred number of workflow runs to return in a page. If not provided, the implementation should use a default page size. The implementation must not return more items than &#x60;page_size&#x60;, but it may return fewer.  Clients should not assume that if fewer than &#x60;page_size&#x60; items are returned that all items have been returned.  The availability of additional pages is indicated by the value of &#x60;next_page_token&#x60; in the response.
        :type page_size: int
        :param page_token: OPTIONAL Token to use to indicate where to start getting results. If unspecified, return the first page of results.
        :type page_token: str

        :rtype: list of dict
        """
        url = self._server_path("query")

        response = requests.request(GET, url)
        res = response.json()["results"]
        self.logger.info("response from remote WES server: %s" % res)

        runs_list = []
        for re in res:
            run_id = re["id"]
            state = self._get_workflow_state_(re["status"], response.status_code)
            runs_list.append(
                RunStatus(
                    run_id=run_id,
                    state=self._build_state_model_(state),
                )
            )

        return RunListResponse(runs=runs_list)

    """In order to execute the workflow, including the primary workflow, tools imported by the workflow, other files 
    referenced by the workflow, or files which are part of the input.

    In the meantime, you should also indict your input files' name with workflowInputs key as json string:
    (for example "{"workflowInputs": "input.json"}") in workflow_params
    """

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
        """Creates a new workflow run and returns a RunId to monitor its progress.

        :return: run_id: str
        """

        # Check if the WES workflow service is healthy. This will throw an exception if not.
        self._check_if_wes_service_healthy_()
        self.logger.debug(f"RUN_WORKFLOW :: wes service is healthy")

        if not workflow_params:
            workflow_params = {}

        # work inside a temp directory context until the request is sent
        # this is needed for any file downloads
        # once complete the temp dir will be removed and all open files will be closed
        with tempfile.TemporaryDirectory() as tmpdir:

            self.logger.debug(f"RUN_WORKFLOW :: tmpdir={tmpdir}")

            data = {
                "workflowType": workflow_type,
                "workflowTypeVersion": workflow_type_version,
            }

            # initialize the files with a placeholder that forces
            # sending a request as multipart/form-data
            # this is requred when there are no files provided in the request
            # and the workflow_url is a remote source
            files = {"labels": (None, json.dumps({"submitted_via": "agc"}))}

            u = urlparse(workflow_url)
            self.logger.debug(f"RUN_WORKFLOW :: urlparse(workflow_url)={u}")
            if u.scheme == "s3":
                self.logger.debug(
                    f"RUN_WORKFLOW :: retrieving '{workflow_url}' => {tmpdir}"
                )
                try:
                    props = get_workflow_from_s3(workflow_url, tmpdir, workflow_type)
                except RuntimeError as e:
                    raise InvalidRequestError(e)
                self.logger.debug(
                    f"RUN_WORKFLOW :: retrieved workflow '{workflow_url}' from S3 => {props}"
                )
                if props.get("data"):
                    data.update(props.get("data"))

                if props.get("files"):
                    files.update(props.get("files"))

            else:
                self.logger.debug(f"RUN_WORKFLOW :: using '{workflow_url}' as is")
                data["workflowUrl"] = workflow_url

            # process workflow attachements
            # this is expected to be only a workflow inputs file
            if workflow_attachment:
                for file in workflow_attachment:
                    self.logger.debug(
                        f"RUN_WORKFLOW :: retrieved workflow attachment : {file.filename}"
                    )

                    if workflow_params.get("workflowInputs"):
                        if file.filename == workflow_params.get("workflowInputs"):
                            # these are inputs supplied at the command line
                            # they take highest priority and should be last on the list
                            if not files.get("workflowInputFiles"):
                                files["workflowInputFiles"] = []

                            files["workflowInputFiles"] += [file]

            # create indexed workflow input keys and files
            # it should be workflowInputs, workflowInputs_2, ... , workflowInputs_5
            if files.get("workflowInputFiles"):
                for i, input_file in enumerate(files["workflowInputFiles"]):
                    j = i + 1

                    if j > 5:
                        raise InvalidRequestError(
                            "maximum number of workflow inputs files exceeded"
                        )

                    ix = "" if j == 1 else f"_{j}"
                    files[f"workflowInputs{ix}"] = input_file

                del files["workflowInputFiles"]

        url = self.url_prefix
        self.logger.info(
            f"RUN_WORKFLOW :: request : url={url}, data={data}, files={files}"
        )

        response = requests.request(POST, url, data=data, files=files)
        self.logger.info(f"RUN_WORKFLOW :: response : {response} ({response.text})")

        if (response.status_code >= 400) and (response.status_code <= 499):
            raise InvalidRequestError(f"HTTP {response.status_code} {response.reason}")
        if response.status_code >= 500 and response.status_code <= 599:
            raise InternalServerError(f"HTTP {response.status_code} {response.reason}")

        return RunId(response.json()["id"])

    def _build_run_log_dict_(self, metadata, outputs):
        run_id = metadata.get("id")
        status = metadata["status"]

        if (status == "fail") or (run_id is None):
            return {
                "run_id": run_id,
                "state": self._get_workflow_state_(status=status),
                "request": {
                    "workflow_params": self.workflow_params,
                    "workflow_type": None,
                    "workflow_type_version": None,
                },
            }

        submitted_files = metadata["submittedFiles"]
        workflow_type = submitted_files["workflowType"]
        workflow_type_version = submitted_files["workflowTypeVersion"]
        workflow_url = submitted_files["workflowUrl"]

        run_request = {
            "workflow_params": self.workflow_params,
            "workflow_type": workflow_type,
            "workflow_type_version": workflow_type_version,
            "workflow_url": workflow_url,
        }

        task_logs = []
        calls = metadata["calls"]

        for task_name in calls.keys():
            for task in calls[task_name]:
                log = {
                    "name": task_name + "|" + task.get("jobId", "XXXXX"),
                    "cmd": [task.get("commandLine")],
                    "start_time": task.get("start"),
                    "end_time": task.get("end"),
                    "stdout": task.get("stdout"),
                    "stderr": task.get("stderr"),
                    "exit_code": ("" if task.get("returnCode") == None else str(task.get("returnCode"))),
                }
                task_logs.append(log)

        run_log_dict = {
            "run_id": run_id,
            "request": run_request,
            "state": self._get_workflow_state_(status=status),
            "task_logs": task_logs,
            "outputs": outputs,
        }
        return run_log_dict

    def _build_run_log_model_(self, run_log_dict):
        request_dict = run_log_dict["request"]
        # Build RunRequest obj
        request = RunRequest(
            workflow_type_version=request_dict["workflow_type_version"],
            workflow_type=request_dict["workflow_type"],
            workflow_params=request_dict["workflow_params"],
        )
        task_logs = []

        for task_log_dict in run_log_dict["task_logs"]:
            # Build Log obj and add it to task_logs list
            task_logs.append(
                Log(
                    name=task_log_dict["name"],
                    cmd=task_log_dict["cmd"],
                    start_time=task_log_dict["start_time"],
                    end_time=task_log_dict["end_time"],
                    stdout=task_log_dict["stdout"],
                    stderr=task_log_dict["stderr"],
                    exit_code=task_log_dict["exit_code"],
                )
            )

        return RunLog(
            run_id=run_log_dict["run_id"],
            request=request,
            state=self._build_state_model_(run_log_dict["state"]),
            task_logs=task_logs,
            outputs=run_log_dict["outputs"],
        )

    # This function just build State model obj from primitive str type
    def _build_state_model_(self, state):
        if state == "EXECUTOR_ERROR":
            return State.EXECUTOR_ERROR
        elif state == "INITIALIZING":
            return State.INITIALIZING
        elif state == "RUNNING":
            return State.RUNNING
        elif state == "COMPLETE":
            return State.COMPLETE
        elif state == "CANCELING":
            return State.CANCELING
        elif state == "CANCELED":
            return State.CANCELED
        elif state == "QUEUED":
            return State.QUEUED
        else:
            return State.UNKNOWN

    """
    This function translate the response of the workflow submitted to the remote WES server to run state defined in Workflow Execution Service. 
    """

    @property
    def workflow_type_versions(self):
        return WORKFLOW_TYPE_VERSIONS

    # TODO: implement SYSTEM_ERROR and QUEUED status
    def _translate_from_response_to_state_(self, response):
        status_code = response.status_code
        status = response.json()["status"]
        return self._get_workflow_state_(status, status_code)

    def _get_workflow_state_(self, status, status_code=None):

        self.logger.info("_get_workflow_state_(%s, %s)" % (status, status_code))
        if (status_code is not None) and (status_code != 200):
            if (status_code >= 400) and (status_code <= 403):
                return "EXECUTOR_ERROR"
            elif status_code == 404:
                return "UNKNOWN"
            else:
                return "SYSTEM_ERROR"
        else:
            if status == "Submitted":
                return "INITIALIZING"
            elif status == "Running":
                return "RUNNING"
            elif status == "Succeeded":
                return "COMPLETE"
            elif status == "Aborting":
                return "CANCELING"
            elif status == "Aborted":
                return "CANCELED"
            elif status == "Failed":
                return "EXECUTOR_ERROR"
            else:
                return "UNKNOWN"

    def _check_if_wes_service_healthy_(self):
        """
        Ping the WES workflow server with a get_service_info call to see if it is healthy.
        Raise an exception if not.
        """

        try:
            response = requests.request(GET, self.health_check_url)
        except:
            self.logger.info(
                "couldn't contact WES service container %s" % self.health_check_url
            )
            raise

        if response.status_code != 200:
            raise InternalServerError

    def _server_path(self, *args):
        args = [str(arg) for arg in args]
        return "/".join([self.url_prefix] + args)


def get_workflow_from_s3(s3_uri: str, localpath: str, workflow_type: str):
    """
    Retrieves a workflow from S3

    :param s3_uri: The S3 URI to the workflow (e.g. s3://bucketname/path/to/workflow.zip)
    :param localpath: The location on the local filesystem to download the workflow
    :param workflow_type: Type of workflow to expect (e.g. wdl, cwl, etc)

    :rtype: dict of `data` and `files`

    If the object is a generic file the file is set as `workflowSource`

    If the object is a `workflow.zip` file containing a single file, that file is set as `workflowSource`

    If the object is a `workflow.zip` file containing multiple files with a MANIFEST.json the MANIFEST is expected to have
      * a mainWorkflowURL property that provides a relative file path in the zip to a workflow file, which will be set as `workflowSource`
      * optionally, if an inputFileURLs property exists that provides a list of relative file paths in the zip to input.json, it will be used to set `workflowInputs`
      * optionally, if an optionFileURL property exists that provides a relative file path in the zip to an options.json file, it will be used to set `workflowOptions`

    If the object is a `workflow.zip` file containing multiple files without a MANIFEST.json
      * a `main` workflow file with an extension matching the workflow_type is expected and will be set as `workflowSource`
      * optionally, if `inputs*.json` files are found in the root level of the zip, they will be set as `workflowInputs(_\d)*` in the order they are found
      * optionally, if an `options.json` file is found in the root level of the zip, it will be set as `workflowOptions`

    If the object is a `workflow.zip` file containing multiple files, the `workflow.zip` file is set as `workflowDependencies`
    """
    s3 = boto3.resource("s3")

    u = urlparse(s3_uri)
    bucket = s3.Bucket(u.netloc)
    key = u.path[1:]

    data = dict()
    files = dict()

    if not key:
        raise RuntimeError("invalid or missing S3 object key")

    try:
        file = path.join(localpath, path.basename(key))
        bucket.download_file(key, file)
    except botocore.exceptions.ClientError as e:
        raise RuntimeError(f"invalid S3 object: {e}")

    if path.basename(file) == "workflow.zip":
        try:
            props = parse_workflow_zip_file(file, workflow_type)
        except Exception as e:
            raise RuntimeError(f"{s3_uri} is not a valid workflow.zip file: {e}")

        if props.get("data"):
            data.update(props.get("data"))

        if props.get("files"):
            files.update(props.get("files"))
    else:
        files["workflowSource"] = open(file, "rb")

    return {"data": data, "files": files}


def parse_workflow_zip_file(file, workflow_type):
    """
    Processes a workflow zip bundle

    :param file: String or Path-like path to a workflow.zip file
    :param workflow_type: String, type of workflow to expect (e.g. "wdl")

    :rtype: dict of `data` and `files`

    If the zip only contains a single file, that file is set as `workflowSource`

    If the zip contains multiple files with a MANIFEST.json file, the MANIFEST is used to determine
    appropriate `data` and `file` arguments. (See: parse_workflow_manifest_file())

    If the zip contains multiple files without a MANIFEST.json file:
      * a `main` workflow file with an extension matching the workflow_type is expected and will be set as `workflowSource`
      * optionally, if `inputs*.json` files are found in the root level of the zip, they will be set as `workflowInputs(_\d)*` in the order they are found
      * optionally, if an `options.json` file is found in the root level of the zip, it will be set as `workflowOptions`

    If the zip contains multiple files, the original zip is set as `workflowDependencies`
    """
    data = dict()
    files = dict()

    wd = path.dirname(file)
    with zipfile.ZipFile(file) as zip:
        zip.extractall(wd)

        contents = zip.namelist()
        if not contents:
            raise RuntimeError("empty workflow.zip")

        if len(contents) == 1:
            # single file workflow
            files["workflowSource"] = open(path.join(wd, contents[0]), "rb")

        else:
            # multifile workflow
            if "MANIFEST.json" in contents:
                props = parse_workflow_manifest_file(path.join(wd, "MANIFEST.json"))

                if props.get("data"):
                    data.update(props.get("data"))

                if props.get("files"):
                    files.update(props.get("files"))

            else:
                if not f"main.{workflow_type.lower()}" in contents:
                    raise RuntimeError(f"'main.{workflow_type}' file not found")

                files["workflowSource"] = open(
                    path.join(wd, f"main.{workflow_type.lower()}"), "rb"
                )

                input_files = [f for f in contents if f.startswith("inputs")]
                if input_files:
                    if not files.get("workflowInputFiles"):
                        files["workflowInputFiles"] = []

                    for input_file in input_files:
                        files[f"workflowInputFiles"] += [
                            open(path.join(wd, input_file), "rb")
                        ]

                if "options.json" in contents:
                    files["workflowOptions"] = open(path.join(wd, "options.json"), "rb")

            # add the original zip bundle as a workflow dependencies file
            files["workflowDependencies"] = open(file, "rb")

    return {"data": data, "files": files}


def parse_workflow_manifest_file(manifest_file):
    """
    Reads a MANIFEST.json file for a workflow zip bundle

    :param manifest_file: String or Path-like path to a MANIFEST.json file

    :rtype: dict of `data` and `files`

    MANIFEST.json is expected to be formatted like:
    .. code-block:: json
       {
           "mainWorkflowURL": "relpath/to/workflow",
           "inputFileURLs": [
               "relpath/to/input-file-1",
               "relpath/to/input-file-2",
               ...
           ],
           "optionsFileURL" "relpath/to/option-file
       }

    The `mainWorkflowURL` property that provides a relative file path in the zip to a workflow file, which will be set as `workflowSource`

    The inputFileURLs property is optional and provides a list of relative file paths in the zip to input.json files. The list is assumed
    to be in the order the inputs should be applied - e.g. higher list index is higher priority. If present, it will be used to set
    `workflowInputs(_\d)` arguments.

    The optionsFileURL property is optional and  provides a relative file path in the zip to an options.json file. If present, it will be
    used to set `workflowOptions`.

    """
    data = dict()
    files = dict()
    with open(manifest_file, "rt") as f:
        manifest = json.loads(f.read())

    u = urlparse(manifest["mainWorkflowURL"])
    if not u.scheme or u.scheme == "file":
        # expect "/path/to/file" or "file:///path/to/file"
        # root is relative to the zip root
        files["workflowSource"] = open(
            workflow_manifest_url_to_path(u, path.dirname(manifest_file)), "rb"
        )

    else:
        data["workflowUrl"] = manifest["mainWorkflowUrl"]

    if manifest.get("inputFileURLs"):
        if not files.get("workflowInputFiles"):
            files["workflowInputFiles"] = []

        for url in manifest["inputFileURLs"]:
            u = urlparse(url)
            if not u.scheme or u.scheme == "file":
                files[f"workflowInputFiles"] += [
                    open(
                        workflow_manifest_url_to_path(u, path.dirname(manifest_file)),
                        "rb",
                    )
                ]

            else:
                raise InvalidRequestError(
                    f"unsupported input file url scheme for: '{url}'"
                )

    if manifest.get("optionsFileURL"):
        u = urlparse(manifest["optionsFileURL"])
        if not u.scheme or u.scheme == "file":
            files["workflowOptions"] = open(
                workflow_manifest_url_to_path(u, path.dirname(manifest_file)), "rb"
            )
        else:
            raise InvalidRequestError(
                f"unsupported option file url scheme for: '{manifest['optionFileURL']}'"
            )

    return {"data": data, "files": files}


def workflow_manifest_url_to_path(url, parent_dir=None):
    relpath = url.path if not url.path.startswith("/") else url.path[1:]
    if parent_dir:
        return path.join(parent_dir, relpath)
    return relpath
