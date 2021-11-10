from unittest.mock import MagicMock

import pytest
from mypy_boto3_batch import BatchClient
from mypy_boto3_logs import CloudWatchLogsClient

from amazon_genomics.wes.adapters.NextflowWESAdapter import NextflowWESAdapter
from rest_api.exception.Exceptions import InternalServerError
from .test_BatchAdapter import generate_batch_job

test_command = ['echo "This is a test!"']

job_queue = "TestJobQueue"
job_definition = "TestJobDefinition"
job_id = "xyz"
job_name = "nextflow"
log_stream = "log-stream"
engine_log_group = "EngineLogGroup"


@pytest.fixture()
def aws_batch() -> BatchClient:
    return MagicMock()


@pytest.fixture()
def aws_logs() -> CloudWatchLogsClient:
    return MagicMock()


@pytest.fixture()
def adapter(aws_batch, aws_logs) -> NextflowWESAdapter:
    return NextflowWESAdapter(
        job_queue=job_queue,
        job_definition=job_definition,
        engine_log_group=engine_log_group,
        aws_batch=aws_batch,
        aws_logs=aws_logs,
    )


def test_get_child_tasks_in_progress(
    aws_batch: BatchClient, aws_logs: CloudWatchLogsClient, adapter: NextflowWESAdapter
):
    job = generate_batch_job({"status": "RUNNING", "startedAt": 1000})

    child_job_id = "child_job_id"
    child_task = generate_batch_job({"jobId": child_job_id, "jobName": "child_task"})

    aws_batch.describe_jobs.side_effect = [
        {"jobs": [child_task]},
    ]

    aws_logs.get_query_results.return_value = {
        "results": [[{"field": "jobId", "value": child_job_id}]],
        "status": "Complete",
    }

    child_tasks = adapter.get_child_tasks(job)
    assert child_tasks == [child_task]


def test_get_child_tasks_completed(
    aws_batch: BatchClient, aws_logs: CloudWatchLogsClient, adapter: NextflowWESAdapter
):
    job = generate_batch_job(
        {"status": "RUNNING", "startedAt": 1000, "stoppedAt": 2000}
    )

    child_job_id = "child_job_id"
    child_task = generate_batch_job({"jobId": child_job_id, "jobName": "child_task"})

    aws_batch.describe_jobs.side_effect = [
        {"jobs": [child_task]},
    ]

    aws_logs.get_query_results.return_value = {
        "results": [[{"field": "jobId", "value": child_job_id}]],
        "status": "Complete",
    }

    child_tasks = adapter.get_child_tasks(job)
    assert child_tasks == [child_task]


def test_get_child_tasks_query_submission_failed(
    aws_logs: CloudWatchLogsClient, adapter: NextflowWESAdapter
):
    job = generate_batch_job(
        {"status": "RUNNING", "startedAt": 1000, "stoppedAt": 2000}
    )

    aws_logs.get_query_results.side_effect = Exception()

    with pytest.raises(Exception):
        adapter.get_child_tasks(job)


def test_get_run_log_query_failed(
    aws_logs: CloudWatchLogsClient, adapter: NextflowWESAdapter
):
    job = generate_batch_job(
        {"status": "RUNNING", "startedAt": 1000, "stoppedAt": 2000}
    )

    aws_logs.get_query_results.return_value = {"results": [], "status": "Failed"}

    with pytest.raises(InternalServerError):
        adapter.get_child_tasks(job)
