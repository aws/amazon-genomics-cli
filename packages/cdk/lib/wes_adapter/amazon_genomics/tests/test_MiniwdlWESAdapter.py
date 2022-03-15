import io
import json
from unittest.mock import MagicMock

import pytest
from botocore.exceptions import ClientError
from botocore.response import StreamingBody
from mypy_boto3_batch import BatchClient
from mypy_boto3_resourcegroupstaggingapi import ResourceGroupsTaggingAPIClient
from mypy_boto3_s3 import S3Client

from amazon_genomics.wes.adapters.MiniWdlWESAdapter import MiniWdlWESAdapter, MINIWDL_OUTPUT_FILE_NAME
from amazon_genomics.wes.adapters.NextflowWESAdapter import NextflowWESAdapter
from .test_BatchAdapter import generate_batch_job

test_command = ['echo "This is a test!"']

job_queue = "TestJobQueue"
job_definition = "TestJobDefinition"
job_id = "xyz"

output_dir_s3_bucket = "output_bucket"
output_dir_s3_prefix = "some/folder"
output_dir_s3_uri = f"s3://{output_dir_s3_bucket}/{output_dir_s3_prefix}"
output_file_path = f"{output_dir_s3_prefix}/{job_id}/{MINIWDL_OUTPUT_FILE_NAME}"


@pytest.fixture()
def aws_batch() -> BatchClient:
    return MagicMock()


@pytest.fixture()
def aws_tags() -> ResourceGroupsTaggingAPIClient:
    return MagicMock()


@pytest.fixture()
def aws_s3() -> S3Client:
    return MagicMock()


@pytest.fixture()
def adapter(aws_batch, aws_tags, aws_s3) -> MiniWdlWESAdapter:
    return MiniWdlWESAdapter(
        job_queue=job_queue,
        job_definition=job_definition,
        output_dir_s3_uri=output_dir_s3_uri,
        aws_batch=aws_batch,
        aws_tags=aws_tags,
        aws_s3=aws_s3
    )


def test_get_child_tasks_in_progress(
    aws_batch: BatchClient,
    aws_tags: ResourceGroupsTaggingAPIClient,
    adapter: MiniWdlWESAdapter,
):
    job = generate_batch_job({"status": "RUNNING", "startedAt": 1000})

    child_job_id = "child_job_id"
    child_task = generate_batch_job({"jobId": child_job_id, "jobName": "child_task"})

    aws_batch.describe_jobs.side_effect = [
        {"jobs": [child_task]},
    ]

    aws_tags.get_resources.return_value = {
        "ResourceTagMappingList": [{"ResourceARN": f"arn:aws:batch:job/{child_job_id}"}]
    }

    child_tasks = adapter.get_child_tasks(job)
    assert child_tasks == [child_task]


def test_get_child_tasks_completed(
    aws_batch: BatchClient,
    aws_tags: ResourceGroupsTaggingAPIClient,
    adapter: NextflowWESAdapter,
):
    job = generate_batch_job(
        {"status": "RUNNING", "startedAt": 1000, "stoppedAt": 2000}
    )

    child_job_id = "child_job_id"
    child_task = generate_batch_job({"jobId": child_job_id, "jobName": "child_task"})

    aws_batch.describe_jobs.side_effect = [
        {"jobs": [child_task]},
    ]

    aws_tags.get_resources.return_value = {
        "ResourceTagMappingList": [{"ResourceARN": f"arn:aws:batch:job/{child_job_id}"}]
    }

    child_tasks = adapter.get_child_tasks(job)
    assert child_tasks == [child_task]


def test_get_child_tasks_query_submission_failed(
    aws_tags: ResourceGroupsTaggingAPIClient, adapter: MiniWdlWESAdapter
):
    job = generate_batch_job(
        {"status": "RUNNING", "startedAt": 1000, "stoppedAt": 2000}
    )

    aws_tags.get_resources.side_effect = Exception()

    with pytest.raises(Exception):
        adapter.get_child_tasks(job)


def test_get_task_output(
    aws_s3: S3Client, adapter: MiniWdlWESAdapter
):
    job = generate_batch_job()
    job_output = {"workflow.output": "somefile.zip"}

    aws_s3.get_object.return_value = mock_s3_object(job_output)

    assert adapter.get_task_outputs(job) == {
        "id": job_id,
        "outputs": job_output
    }


def test_get_task_output_no_file(
    aws_s3: S3Client, adapter: MiniWdlWESAdapter
):
    job = generate_batch_job()
    aws_s3.get_object.side_effect = ClientError(
        error_response={
            'Error': {'Code': "NoSuchKey"}
        },
        operation_name="GetObject"
    )

    assert adapter.get_task_outputs(job) == {
        "id": job_id,
        "outputs": None
    }


def test_get_task_output_exception(
    aws_s3: S3Client, adapter: MiniWdlWESAdapter
):
    job = generate_batch_job()
    aws_s3.get_object.side_effect = Exception()

    with pytest.raises(Exception):
        adapter.get_task_outputs(job)


def mock_s3_object(obj):
    body_encoded = json.dumps(obj).encode()
    body = StreamingBody(
        io.BytesIO(body_encoded),
        len(body_encoded)
    )
    return {
        "Body": body
    }
