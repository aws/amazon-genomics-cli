from typing import Dict, List

from rest_api.models import (
    ServiceInfo,
    RunId,
    RunListResponse,
    RunStatus,
    RunLog,
    RunRequest,
)


def verify_service_info(service_info):
    assert isinstance(service_info, ServiceInfo)

    # ensure all required fields exist and are correct type
    assert service_info.supported_wes_versions and isinstance(
        service_info.supported_wes_versions, List
    )
    assert service_info.workflow_type_versions and isinstance(
        service_info.workflow_type_versions, Dict
    )


def verify_run_id(runId):
    assert runId
    assert isinstance(runId, RunId)


def verify_run_list_response(res):
    assert res
    assert isinstance(res, RunListResponse)


def verify_run_status(status):
    assert status
    assert isinstance(status, RunStatus)


def verify_run_log(log):
    assert log
    assert isinstance(log, RunLog)
    assert log.run_id
    assert log.request and isinstance(log.request, RunRequest)
