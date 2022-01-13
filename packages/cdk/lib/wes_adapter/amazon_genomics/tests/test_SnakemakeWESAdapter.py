from unittest.mock import MagicMock

import pytest
from mypy_boto3_batch import BatchClient
from mypy_boto3_logs import CloudWatchLogsClient

from amazon_genomics.wes.adapters.SnakemakeWESAdapter import SnakemakeWESAdapter
from rest_api.exception.Exceptions import InternalServerError
from .test_BatchAdapter import generate_batch_job

test_command = ['echo "This is a test!"']

job_queue = "TestJobQueue"
job_definition = "TestJobDefinition"
job_id = "xyz"
job_name = "snakemake"
log_stream = "log-stream"
engine_log_group = "EngineLogGroup"


@pytest.fixture()
def aws_batch() -> BatchClient:
    return MagicMock()


@pytest.fixture()
def aws_logs() -> CloudWatchLogsClient:
    return MagicMock()


@pytest.fixture()
def adapter(aws_batch, aws_logs) -> SnakemakeWESAdapter:
    return SnakemakeWESAdapter(
        job_queue=job_queue,
        job_definition=job_definition,
        engine_log_group=engine_log_group,
        aws_batch=aws_batch,
        aws_logs=aws_logs,
    )


def test_get_child_tasks_in_progress(
    aws_batch: BatchClient, aws_logs: CloudWatchLogsClient, adapter: SnakemakeWESAdapter
):
    job = generate_batch_job({"status": "RUNNING", "startedAt": 1000})
    child_tasks = adapter.get_child_tasks(job)
    assert child_tasks == []
