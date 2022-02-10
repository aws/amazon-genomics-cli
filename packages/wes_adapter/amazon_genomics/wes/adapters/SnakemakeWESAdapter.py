import typing
import json


import os

import boto3
from botocore.exceptions import ClientError
from amazon_genomics.wes.adapters.util.util import describe_batch_jobs_with_tag
from mypy_boto3_batch import BatchClient
from mypy_boto3_batch.type_defs import JobDetailTypeDef
from mypy_boto3_logs.type_defs import ResultFieldTypeDef
from mypy_boto3_resourcegroupstaggingapi import ResourceGroupsTaggingAPIClient
from mypy_boto3_s3 import S3Client
from amazon_genomics.wes.adapters.BatchAdapter import BatchAdapter
from rest_api.models import (
    WorkflowTypeVersion,
    ServiceInfo,
)


SM_PARENT_TAG_KEY = "AWS_BATCH_PARENT_JOB_ID"
SM_OUTPUT_FILE_NAME = "outputs.json"

TASK_QUEUE = os.getenv("TASK_QUEUE")
WORKFLOW_ROLE = os.getenv("WORKFLOW_ROLE")
FSAP_ID = os.getenv("FSAP_ID")


class SnakemakeWESAdapter(BatchAdapter):
    """
    WES controller that handles WES requests for a Snakemake workflow engine
    """

    def __init__(
        self,
        job_queue: str,
        job_definition: str,
        output_dir_s3_uri: str,
        aws_batch: BatchClient = None,
        aws_tags: ResourceGroupsTaggingAPIClient = None,
        aws_s3: S3Client = None,
        logger=None,
    ):
        super().__init__(job_queue, job_definition, aws_batch, logger)
        self.output_dir_s3_uri = output_dir_s3_uri
        self.task_queue = TASK_QUEUE
        self.workflow_role = WORKFLOW_ROLE
        self.fsap_id = FSAP_ID
        self.aws_tags: ResourceGroupsTaggingAPIClient = aws_tags or boto3.client(
            "resourcegroupstaggingapi", region_name=os.environ["AWS_REGION"]
        )
        self.aws_s3: S3Client = aws_s3 or boto3.client(
            "s3", region_name=os.environ["AWS_REGION"]
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
        engine_params_to_pass = []
        if workflow_engine_parameters is not None:
            print("The engine parameters are {}".format(workflow_engine_parameters))
            engine_params_to_pass.append(workflow_engine_parameters)

        engine_params_to_pass.extend(
            [
                "--aws-batch",
                "--cores all",
                "--aws-batch-workflow-role {}".format(self.workflow_role),
                "--aws-batch-task-queue {}".format(self.task_queue),
                "--aws-batch-fsap-id {}".format(self.fsap_id),
            ]
        )
        delimiter = " "
        command = [workflow_url, delimiter.join(engine_params_to_pass)]
        return command

    def environment(self):
        return {}

    def get_child_tasks(
        self, head_job: JobDetailTypeDef
    ) -> typing.List[JobDetailTypeDef]:
        return describe_batch_jobs_with_tag(
            tag_key=SM_PARENT_TAG_KEY,
            tag_value=head_job["jobId"],
            aws_batch=self.aws_batch,
            aws_tags=self.aws_tags,
        )

    def get_task_outputs(self, head_job: JobDetailTypeDef):
        # TODO: update implementation based on executor s3 write changes
        job_id = head_job.get("jobId")
        bucket, folder = self.output_dir_s3_uri.split("/", 2)[-1].split("/", 1)
        output_file_key = f"{folder}/{job_id}/{SM_OUTPUT_FILE_NAME}"
        output = self.get_s3_object_json(bucket=bucket, output_file_key=output_file_key)
        return {"id": job_id, "outputs": output}

    def get_s3_object_json(self, bucket, output_file_key):
        try:
            output_object = self.aws_s3.get_object(Bucket=bucket, Key=output_file_key)
            return json.load(output_object["Body"])
        except ClientError as ex:
            if ex.response["Error"]["Code"] == "NoSuchKey":
                self.logger.warn(f"No object found")
                return None
            else:
                raise ex

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
