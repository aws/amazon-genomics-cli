import traceback
import os
import typing
import time
from abc import abstractmethod
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime
from typing import Optional
from typing import Iterable

import boto3
from mypy_boto3_batch import BatchClient
from mypy_boto3_batch.type_defs import JobDetailTypeDef

from amazon_genomics.wes.adapters import AbstractWESAdapter
from rest_api.exception.Exceptions import InternalServerError
from rest_api.models import (
    Log,
    RunLog,
    RunListResponse,
    State,
    RunId,
    RunStatus,
)

USER_CANCELLATION_REASON = "User Canceled"


class BatchAdapter(AbstractWESAdapter):
    """
    Base class for Adapters which submit jobs to AWS Batch
    """

    def __init__(
            self,
            job_queue: str,
            job_definition: str,
            aws_batch: BatchClient = None,
            logger=None,
    ):
        super().__init__(logger)
        self.job_queue = job_queue
        self.job_definition = job_definition
        self.aws_batch: BatchClient = (
            aws_batch
            if aws_batch
            else boto3.client("batch", region_name=os.environ["AWS_REGION"])
        )

    @abstractmethod
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
        pass

    @abstractmethod
    def environment(self):
        pass

    def cancel_run(self, run_id: str) -> Optional[RunId]:
        """Cancel a running workflow.

        :param run_id:
        :type run_id: str

        :rtype: RunId
        """
        try:
            self.aws_batch.terminate_job(jobId=run_id, reason=USER_CANCELLATION_REASON)
            return RunId(run_id)
        except Exception as e:
            traceback.print_exc()
            raise InternalServerError(f"Failed to cancel task for {run_id}", e)

    def get_run_log(self, run_id) -> Optional[RunLog]:
        head_job = self.describe_job(run_id)
        if not head_job:
            return None

        try:
            child_jobs = self.get_child_tasks(head_job)
            task_logs = list(map(lambda task_job: self.to_log(task_job), child_jobs))
        except Exception as e:
            traceback.print_exc()
            raise InternalServerError(f"Failed to load child tasks for job {run_id}", e)

        return RunLog(
            run_id=run_id,
            state=self.batch_job_wes_state(
                head_job["status"], head_job.get("statusReason", "")
            ),
            run_log=self.to_log(head_job),
            task_logs=task_logs,
            outputs=self.get_task_outputs(head_job),
        )

    @abstractmethod
    def get_task_outputs(self, task: JobDetailTypeDef):
        pass

    def get_run_status(self, run_id) -> RunStatus:
        job = self.describe_job(run_id)
        if not job:
            return None
        return self.to_run_status(
            job_id=job["jobId"],
            job_status=job["status"],
            job_status_reason=job.get("statusReason", ""),
        )

    def list_runs(self, page_size=None, page_token=None) -> RunListResponse:
        list_jobs_response = self.aws_batch.list_jobs(
            jobQueue=self.job_queue,
            maxResults=page_size or 50,
            nextToken=page_token or "",
            filters=[
                {"name": "JOB_DEFINITION", "values": [self.job_definition]},
            ],
        )

        next_token = list_jobs_response.get("nextToken")
        job_summaries = list_jobs_response["jobSummaryList"]

        runs = list(
            map(
                lambda job_summary: self.to_run_status(
                    job_id=job_summary["jobId"],
                    job_status=job_summary["status"],
                    job_status_reason=job_summary.get("statusReason", ""),
                ),
                job_summaries,
            )
        )
        return RunListResponse(runs=runs, next_page_token=next_token)

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
        """
        Submit "workflow job" based on given configuration details; return the Batch job uuid
        """
        command = self.command(
            workflow_params=workflow_params,
            workflow_type=workflow_type,
            workflow_type_version=workflow_type_version,
            tags=tags,
            workflow_engine_parameters=workflow_engine_parameters,
            workflow_url=workflow_url,
            workflow_attachment=workflow_attachment,
        )

        submit_job_response = self.aws_batch.submit_job(
            jobName="agc-run-workflow",
            jobQueue=self.job_queue,
            jobDefinition=self.job_definition,
            containerOverrides={
                "command": command,
            },
        )
        return RunId(submit_job_response["jobId"])

    def describe_job(self, job_id: str) -> Optional[JobDetailTypeDef]:
        jobs = self.describe_jobs([job_id])
        if not jobs:
            return None
        else:
            return jobs[0]

    def describe_jobs(self, job_ids: typing.List[str]) -> typing.List[JobDetailTypeDef]:
        if not job_ids:
            return []

        jobs = []
        job_ids_sets = chunks(job_ids, 100)
        with ThreadPoolExecutor(max_workers=10) as executor:
            future_jobs = {
                executor.submit(
                    self.aws_batch.describe_jobs, jobs=job_ids_set
                ): job_ids_set
                for job_ids_set in job_ids_sets
            }
            for future in as_completed(future_jobs):
                job_ids_set = future_jobs[future]
                try:
                    response = future.result(timeout=5)
                    jobs += response["jobs"]
                except Exception as e:
                    self.logger.error(f"error retrieving jobs: {e}")

        return jobs

    @abstractmethod
    def get_child_tasks(
            self, head_job: JobDetailTypeDef
    ) -> typing.List[JobDetailTypeDef]:
        pass

    @staticmethod
    def batch_job_wes_state(job_status, job_status_reason) -> State:
        """
        Derive WES job state from AWS Batch job description (of workflow job)
        :param job_status: Job status as given by AWS Batch.
            See: https://docs.aws.amazon.com/batch/latest/userguide/job_states.html
        :param job_status_reason: Status reason used to distinguish CANCELED from FAILED run status for WES
            See: https://docs.aws.amazon.com/batch/latest/APIReference/API_JobSummary.html
        :return: WES Job status equivalent.
            See: https://ga4gh.github.io/workflow-execution-service-schemas/docs/#operation/GetRunStatus
        """
        #
        if job_status in ("SUBMITTED", "PENDING", "RUNNABLE"):
            return State.QUEUED
        elif job_status == "STARTING":
            return State.INITIALIZING
        elif job_status == "RUNNING":
            return State.RUNNING
        elif job_status == "SUCCEEDED":
            return State.COMPLETE
        elif job_status == "FAILED":
            if job_status_reason == USER_CANCELLATION_REASON:
                return State.CANCELED
            # TODO: detect SYSTEM_ERROR (probably, anything but "Essential container in task exited")
            return State.EXECUTOR_ERROR
        else:
            return State.UNKNOWN

    @staticmethod
    def to_run_status(job_id: str, job_status: str, job_status_reason: str):
        return RunStatus(
            run_id=job_id,
            state=BatchAdapter.batch_job_wes_state(job_status, job_status_reason),
        )

    @staticmethod
    def to_log(job_details: JobDetailTypeDef) -> Log:
        # The CLI expects a pipe-delimited string of name + string in the response
        task_name = f"{job_details['jobName']}|{job_details['jobId']}"

        start_time = to_iso(job_details.get("startedAt"))
        end_time = to_iso(job_details.get("stoppedAt"))

        return Log(
            name=task_name,
            cmd=job_details["container"]["command"],
            start_time=start_time,
            end_time=end_time,
            stdout=job_details["container"].get("logStreamName"),
            exit_code=job_details["container"].get("exitCode"),
        )


def to_iso(epoch: Optional[int]) -> Optional[str]:
    if not epoch:
        return None
    return datetime.utcfromtimestamp(epoch / 1000.0).astimezone().isoformat()


def chunks(l: list, n: int) -> Iterable[list]:
    """split list l into chunks of size n"""
    for i in range(0, len(l), n):
        yield l[i: i + n]
