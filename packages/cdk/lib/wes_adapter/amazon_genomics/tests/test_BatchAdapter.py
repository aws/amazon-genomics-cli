import typing
from unittest import mock
from unittest.mock import MagicMock

import pytest
from mypy_boto3_batch import BatchClient
from mypy_boto3_batch.type_defs import JobDetailTypeDef
from amazon_genomics.wes.adapters.BatchAdapter import BatchAdapter
from rest_api.exception.Exceptions import InternalServerError
from rest_api.models import (
    State,
    Log,
    RunLog,
    RunListResponse,
    RunStatus,
    RunId,
    ServiceInfo,
    WorkflowTypeVersion,
)

test_command = ['echo "This is a test!"']

job_queue = "TestJobQueue"
job_definition = "TestJobDefinition"
job_id = "xyz"
job_name = "agc-run-workflow"
log_stream = "log-stream"
engine_log_group = "EngineLogGroup"


class StubBatchAdapter(BatchAdapter):
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
        return [workflow_url]

    def environment(self):
        return {}

    def get_child_tasks(
        self, head_job: JobDetailTypeDef
    ) -> typing.List[JobDetailTypeDef]:
        return []

    @property
    def workflow_type_versions(self):
        return {"STUBENGINE": WorkflowTypeVersion(["1.0", "dsl2"])}

    def get_task_outputs(self, head_job: JobDetailTypeDef):
        return {
            "id": head_job.get("jobId"),
        }

    def get_service_info(self):
        return ServiceInfo(
            supported_wes_versions=self.supported_wes_versions,
            workflow_type_versions=self.workflow_type_versions,
        )


@pytest.fixture
def aws_batch() -> BatchClient:
    return MagicMock()


@pytest.fixture
def adapter(aws_batch) -> StubBatchAdapter:
    return StubBatchAdapter(
        job_queue=job_queue, job_definition=job_definition, aws_batch=aws_batch
    )


def test_list_runs_no_runs(aws_batch: BatchClient, adapter: StubBatchAdapter):
    aws_batch.list_jobs.return_value = {"jobSummaryList": []}
    list_runs_response = adapter.list_runs()
    assert not list_runs_response.runs
    assert list_runs_response.next_page_token is None


def test_list_runs_single_page(aws_batch: BatchClient, adapter: StubBatchAdapter):
    aws_batch.list_jobs.return_value = {
        "jobSummaryList": [
            {"jobId": "abc", "status": "RUNNING"},
            {"jobId": job_id, "status": "SUCCEEDED"},
        ]
    }
    list_runs_response = adapter.list_runs()
    expected_response = RunListResponse(
        runs=[
            RunStatus(run_id="abc", state=State.RUNNING),
            RunStatus(run_id=job_id, state=State.COMPLETE),
        ]
    )
    assert list_runs_response == expected_response


def test_list_runs_multiple_pages(aws_batch: BatchClient, adapter: StubBatchAdapter):
    aws_batch.list_jobs.return_value = {
        "jobSummaryList": [
            {"jobId": "abc", "status": "RUNNING"},
            {"jobId": job_id, "status": "SUCCEEDED"},
        ],
        "nextToken": "someToken",
    }
    list_runs_response = adapter.list_runs()
    expected_response = RunListResponse(
        runs=[
            RunStatus(run_id="abc", state=State.RUNNING),
            RunStatus(run_id=job_id, state=State.COMPLETE),
        ],
        next_page_token="someToken",
    )
    assert list_runs_response == expected_response


def test_get_run_status_no_runs(aws_batch: BatchClient, adapter: StubBatchAdapter):
    aws_batch.describe_jobs.return_value = {"jobs": []}
    run_status = adapter.get_run_status(job_id)
    assert run_status is None


def test_get_run_status_with_valid_response(
    aws_batch: BatchClient, adapter: StubBatchAdapter
):
    aws_batch.describe_jobs.return_value = {
        "jobs": [generate_batch_job({"status": "SUCCEEDED"})]
    }
    run_status = adapter.get_run_status(job_id)
    expected_status = RunStatus(
        run_id=job_id,
        state=State.COMPLETE,
    )
    assert run_status == expected_status


def test_get_run_status_with_cancelled_job(
    aws_batch: BatchClient, adapter: StubBatchAdapter
):
    aws_batch.describe_jobs.return_value = {
        "jobs": [
            generate_batch_job({"status": "FAILED", "statusReason": "User Canceled"})
        ]
    }
    run_status = adapter.get_run_status(job_id)
    expected_status = RunStatus(
        run_id=job_id,
        state=State.CANCELED,
    )
    assert run_status == expected_status


def test_get_run_log_nonexistent_job(aws_batch: BatchClient, adapter: StubBatchAdapter):
    aws_batch.describe_jobs.return_value = {"jobs": []}
    run_log = adapter.get_run_log(job_id)
    assert run_log is None


def test_get_run_log_not_started(aws_batch: BatchClient, adapter: StubBatchAdapter):
    job = generate_batch_job({"status": "RUNNABLE"})

    aws_batch.describe_jobs.return_value = {"jobs": [job]}
    run_log = adapter.get_run_log(job_id)
    assert run_log == RunLog(
        run_id=job_id,
        state=State.QUEUED,
        run_log=Log(name="agc-run-workflow|xyz", cmd=test_command, stdout=log_stream),
        task_logs=[],
        outputs={"id": job_id},
    )


def test_get_run_log_in_progress(aws_batch: BatchClient, adapter: StubBatchAdapter):
    job = generate_batch_job({"status": "RUNNING", "startedAt": 1000})

    child_job_id = "child_job_id"
    child_task = generate_batch_job({"jobId": child_job_id, "jobName": "child_task"})

    aws_batch.describe_jobs.side_effect = [
        {"jobs": [job]},
        {"jobs": [child_task]},
    ]

    with mock.patch.object(adapter, "get_child_tasks", return_value=[child_task]):
        run_log = adapter.get_run_log(job_id)
        assert run_log == RunLog(
            run_id=job_id,
            state=State.RUNNING,
            run_log=Log(
                name="agc-run-workflow|xyz",
                cmd=test_command,
                stdout=log_stream,
                start_time="1970-01-01T00:00:01+00:00",
            ),
            task_logs=[
                Log(
                    name="child_task|child_job_id",
                    cmd=test_command,
                    stdout=log_stream,
                )
            ],
            outputs={"id": job_id},
        )


def test_get_run_log_completed(aws_batch: BatchClient, adapter: StubBatchAdapter):
    job = generate_batch_job(
        {"status": "RUNNING", "startedAt": 1000, "stoppedAt": 2000}
    )

    child_job_id = "child_job_id"
    child_task = generate_batch_job({"jobId": child_job_id, "jobName": "child_task"})

    aws_batch.describe_jobs.side_effect = [
        {"jobs": [job]},
        {"jobs": [child_task]},
    ]

    with mock.patch.object(adapter, "get_child_tasks", return_value=[child_task]):
        run_log = adapter.get_run_log(job_id)
        assert run_log == RunLog(
            run_id=job_id,
            state=State.RUNNING,
            run_log=Log(
                name="agc-run-workflow|xyz",
                cmd=test_command,
                stdout=log_stream,
                start_time="1970-01-01T00:00:01+00:00",
                end_time="1970-01-01T00:00:02+00:00",
            ),
            task_logs=[
                Log(
                    name="child_task|child_job_id",
                    cmd=test_command,
                    stdout=log_stream,
                )
            ],
            outputs={"id": job_id},
        )


def test_run_workflow(aws_batch: BatchClient, adapter: StubBatchAdapter):
    workflow_url = "s3://my_workflow/"
    aws_batch.submit_job.return_value = {"jobId": job_id}
    batch_job_id = adapter.run_workflow(workflow_url=workflow_url)
    aws_batch.submit_job.assert_called_with(
        jobName=job_name,
        jobQueue=job_queue,
        jobDefinition=job_definition,
        containerOverrides={
            "command": [workflow_url],
        },
    )
    assert batch_job_id.run_id == job_id


def test_get_service_info(adapter: StubBatchAdapter):
    service_info = adapter.get_service_info()
    assert service_info.supported_wes_versions == ["1.0.0"]
    assert service_info.workflow_type_versions["STUBENGINE"].workflow_type_version == [
        "1.0",
        "dsl2",
    ]


def test_cancel_run(aws_batch: BatchClient, adapter: StubBatchAdapter):
    aws_batch.terminate_job.return_value = job_id
    canceled_job_id = adapter.cancel_run(job_id)
    assert canceled_job_id == RunId(job_id)


def test_cancel_run_failed(aws_batch: BatchClient, adapter: StubBatchAdapter):
    aws_batch.terminate_job.side_effect = Exception()
    with pytest.raises(InternalServerError):
        adapter.cancel_run(job_id)


def generate_batch_job(overrides=None):
    job_defaults = {
        "jobId": job_id,
        "jobName": job_name,
        "status": "RUNNABLE",
        "container": {"command": test_command, "logStreamName": log_stream},
    }
    return {**job_defaults, **overrides} if overrides else job_defaults
