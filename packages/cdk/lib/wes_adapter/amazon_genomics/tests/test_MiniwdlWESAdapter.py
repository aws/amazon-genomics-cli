from unittest.mock import MagicMock

import pytest
from mypy_boto3_batch import BatchClient
from mypy_boto3_resourcegroupstaggingapi import ResourceGroupsTaggingAPIClient
from amazon_genomics.wes.adapters.NextflowWESAdapter import NextflowWESAdapter
from amazon_genomics.wes.adapters.MiniWdlWESAdapter import MiniWdlWESAdapter
from .test_BatchAdapter import generate_batch_job

test_command = ['echo "This is a test!"']

job_queue = "TestJobQueue"
job_definition = "TestJobDefinition"
job_id = "xyz"


@pytest.fixture()
def aws_batch() -> BatchClient:
    return MagicMock()


@pytest.fixture()
def aws_tags() -> ResourceGroupsTaggingAPIClient:
    return MagicMock()


@pytest.fixture()
def adapter(aws_batch, aws_tags) -> MiniWdlWESAdapter:
    return MiniWdlWESAdapter(
        job_queue=job_queue,
        job_definition=job_definition,
        aws_batch=aws_batch,
        aws_tags=aws_tags,
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
    aws_tags: ResourceGroupsTaggingAPIClient, adapter: NextflowWESAdapter
):
    job = generate_batch_job(
        {"status": "RUNNING", "startedAt": 1000, "stoppedAt": 2000}
    )

    aws_tags.get_resources.side_effect = Exception()

    with pytest.raises(Exception):
        adapter.get_child_tasks(job)
