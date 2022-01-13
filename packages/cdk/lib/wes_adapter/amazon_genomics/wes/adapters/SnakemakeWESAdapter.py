import typing

import time
import os

import boto3
from mypy_boto3_batch import BatchClient
from mypy_boto3_batch.type_defs import JobDetailTypeDef
from mypy_boto3_logs import CloudWatchLogsClient
from mypy_boto3_logs.type_defs import ResultFieldTypeDef

from amazon_genomics.wes.adapters.BatchAdapter import BatchAdapter
from rest_api.models import (
    WorkflowTypeVersion,
    ServiceInfo,
)


class SnakemakeWESAdapter(BatchAdapter):
    """
    WES controller that handles WES requests for a Snakemake workflow engine
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
        command = [workflow_url]
        return command

    def environment(self):
        return {}

    def get_child_tasks(
        self, head_job: JobDetailTypeDef
    ) -> typing.List[JobDetailTypeDef]:
        # TODO: implement finding child tasks for snakemake
        return []

    def get_task_outputs(self, task: JobDetailTypeDef):
        # TODO: implement finding task outputs for snakemake
        return {}

    @property
    def workflow_type_versions(self):
        # TODO: update to supported snakemake version after off internal-fork
        return {"SNAKEMAKE": WorkflowTypeVersion(["1.0"])}

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
