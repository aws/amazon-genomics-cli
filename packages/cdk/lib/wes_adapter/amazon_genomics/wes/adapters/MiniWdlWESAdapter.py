import json
import os
import typing

import boto3
from botocore.exceptions import ClientError
from mypy_boto3_batch import BatchClient
from mypy_boto3_batch.type_defs import JobDetailTypeDef
from mypy_boto3_resourcegroupstaggingapi import ResourceGroupsTaggingAPIClient
from mypy_boto3_s3 import S3Client
from amazon_genomics.wes.adapters.BatchAdapter import BatchAdapter
from rest_api.models import (
    WorkflowTypeVersion,
    ServiceInfo,
)

MINIWDL_PARENT_TAG_KEY = "AWS_BATCH_PARENT_JOB_ID"
MINIWDL_OUTPUT_FILE_NAME = "outputs.json"

class MiniWdlWESAdapter(BatchAdapter):
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
        command = [workflow_url]
        return command

    def environment(self):
        return {}

    def get_child_tasks(
        self, head_job: JobDetailTypeDef
    ) -> typing.List[JobDetailTypeDef]:
        return self.describe_batch_jobs_with_tag(
            tag_key=MINIWDL_PARENT_TAG_KEY, tag_value=head_job["jobId"]
        )

    @property
    def workflow_type_versions(self):
        return {"WDL": WorkflowTypeVersion(["1.0", "draft-2"])}

    def describe_batch_jobs_with_tag(
        self,
        tag_key,
        tag_value,
    ):
        """
        Retrieve descriptions of all Batch jobs with the given tag
        """
        pagination_token = None
        all_descriptions = []
        get_resources_kwargs = {
            "TagFilters": [{"Key": tag_key, "Values": [tag_value]}],
            "ResourceTypeFilters": ["batch:job"],
        }
        while True:
            if pagination_token:
                get_resources_kwargs["PaginationToken"] = pagination_token
            resources = self.aws_tags.get_resources(**get_resources_kwargs)
            resource_tag_mappings = resources.get("ResourceTagMappingList", [])
            job_arns = map(
                lambda tag_mapping: tag_mapping["ResourceARN"], resource_tag_mappings
            )
            job_ids = list(map(job_id_from_arn, job_arns))
            if job_ids:
                descriptions = self.aws_batch.describe_jobs(jobs=job_ids)["jobs"]
                all_descriptions += descriptions
            pagination_token = resources.get("PaginationToken", None)
            if not pagination_token:
                return all_descriptions

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

    def get_task_outputs(self, head_job: JobDetailTypeDef):
        job_id = head_job.get("jobId")
        bucket, folder = self.output_dir_s3_uri.split('/', 2)[-1].split('/', 1)
        output_file_key = f"{folder}/{job_id}/{MINIWDL_OUTPUT_FILE_NAME}"
        output = self.get_s3_object_json(bucket, output_file_key)
        return {
            "id": job_id,
            "outputs": output
        }

    def get_s3_object_json(self, bucket, output_file_key):
        try:
            output_object = self.aws_s3.get_object(
                Bucket=bucket,
                Key=output_file_key
            )
            return json.load(output_object["Body"])
        except ClientError as ex:
            if ex.response['Error']['Code'] == 'NoSuchKey':
                self.logger.warn(f"No object found")
                return None
            else:
                raise ex

def job_id_from_arn(job_arn: str) -> str:
    return job_arn[job_arn.rindex("/") + 1 :]
