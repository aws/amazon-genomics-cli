import typing

import time
import os

import boto3
from mypy_boto3_batch import BatchClient
from mypy_boto3_batch.type_defs import JobDetailTypeDef
from mypy_boto3_logs import CloudWatchLogsClient
from mypy_boto3_logs.type_defs import ResultFieldTypeDef

from amazon_genomics.wes.adapters.BatchAdapter import BatchAdapter
from rest_api.exception.Exceptions import InternalServerError
from rest_api.models import (
    WorkflowTypeVersion,
    ServiceInfo,
)


class NextflowWESAdapter(BatchAdapter):
    """
    WES controller that handles WES requests for a Nextflow workflow engine
    """

    def __init__(
        self,
        job_queue: str,
        job_definition: str,
        engine_log_group: str,
        aws_batch: BatchClient = None,
        aws_logs: CloudWatchLogsClient = None,
        logger=None,
    ):
        super().__init__(job_queue, job_definition, aws_batch, logger)
        self.engine_log_group = engine_log_group
        self.aws_logs: CloudWatchLogsClient = aws_logs or boto3.client(
            "logs", region_name=os.environ["AWS_REGION"]
        )

    def command(
        self,
        workflow_params=None,
        workflow_type=None,
        workflow_type_version=None,
        tags=None,
        workflow_engine_parameters=None,
        workflow_url=None,
        workflow_attachment=None,
    ):
        command = [
            workflow_url,
        ]
        """
        Support DSL 2 if specified
        https://www.nextflow.io/docs/latest/dsl2.html
        """
        if workflow_type_version == "dsl2":
            command.append("-dsl2")
        """
        TODO: Add support for params-file.
        """
        return command

    def environment(self):
        return {}

    def get_task_outputs(self, head_job: JobDetailTypeDef):
        outputs = []
        if "logStreamName" in head_job["container"]:
            log_stream = head_job["container"]["logStreamName"]
            query_string = f"""
                fields @message, @logStream
                | filter @logStream = "{log_stream}"
                | filter @message like /TaskPollingMonitor - Task completed/
                | parse 'name: *;' as name 
                | parse 'id: *;' as id 
                | parse 'status: *;' as status
                | parse 'exit: *;' as exit
                | parse 'error: *;' as error
                | parse 'workDir: *]' as workDir
                | display id, name,status, exit, error, workDir
            """
            outputs = self.query_logs_for_job(head_job, query_string)
        return {
            "id": head_job.get("jobId"),
            "outputs": outputs,
        }

    def get_child_tasks(
        self, head_job: JobDetailTypeDef
    ) -> typing.List[JobDetailTypeDef]:
        if "logStreamName" not in head_job["container"]:
            return []
        log_stream = head_job["container"]["logStreamName"]
        query_string = f"""
        fields @message, @logStream
        | filter @logStream = "{log_stream}"
        | filter @message like /\[AWS BATCH\] submitted/
        | parse 'job=*;' as jobId
        | stats latest(@ingestionTime) by jobId
        | display jobId
        """
        jobs = self.query_logs_for_job(head_job, query_string)
        child_job_ids = list(map(lambda job: job["jobId"], jobs))
        return self.describe_jobs(child_job_ids)

    @property
    def workflow_type_versions(self):
        return {"NEXTFLOW": WorkflowTypeVersion(["1.0", "dsl2"])}

    def get_service_info(self):
        """Get information about Workflow Execution Service.

        May include information related (but not limited to) the workflow
        descriptor formats, versions supported, the WES API versions supported,
        and information about general service availability.

        :rtype: ServiceInfo
        """
        return ServiceInfo(
            supported_wes_versions=self.supported_wes_versions,
            workflow_type_versions=self.workflow_type_versions,
        )

    def query_logs_for_job(self, job_details: JobDetailTypeDef, query: str):
        start_time = job_details.get("startedAt")
        if not start_time:
            # the job is not started yet, so no tasks will have been created.
            return []

        end_time = job_details.get("stoppedAt")

        query_id = self.aws_logs.start_query(
            logGroupName=self.engine_log_group,
            startTime=start_time,
            endTime=end_time or int(time.time()),
            queryString=query,
            # TODO: handle pagination? GetRunLog doesn't seem to support it...
            limit=100,
        )["queryId"]
        response = None

        while response is None or response["status"] in ("Scheduled", "Running"):
            self.logger.info(f"Waiting for query [{query_id}] to complete ...")
            time.sleep(1)
            response = self.aws_logs.get_query_results(queryId=query_id)
        if response["status"] != "Complete":
            raise InternalServerError("Logs query for child tasks was not successful")

        results = list(map(lambda result: to_dict(result), response["results"]))
        return results


# Cloudwatch Log Query results are a list of 'field' and 'value' named tuples
# This function converts them into a dict
def to_dict(results: typing.List[ResultFieldTypeDef]):
    return {
        result["field"]: result["value"]
        for result in results
        if result["field"] != "@ptr"
    }
