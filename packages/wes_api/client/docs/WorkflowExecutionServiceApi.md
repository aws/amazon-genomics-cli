# \WorkflowExecutionServiceApi

All URIs are relative to *http://localhost/ga4gh/wes/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CancelRun**](WorkflowExecutionServiceApi.md#CancelRun) | **Post** /runs/{run_id}/cancel | Cancel a running workflow.
[**GetRunLog**](WorkflowExecutionServiceApi.md#GetRunLog) | **Get** /runs/{run_id} | Get detailed info about a workflow run.
[**GetRunStatus**](WorkflowExecutionServiceApi.md#GetRunStatus) | **Get** /runs/{run_id}/status | Get quick status info about a workflow run.
[**GetServiceInfo**](WorkflowExecutionServiceApi.md#GetServiceInfo) | **Get** /service-info | Get information about Workflow Execution Service.
[**ListRuns**](WorkflowExecutionServiceApi.md#ListRuns) | **Get** /runs | List the workflow runs.
[**RunWorkflow**](WorkflowExecutionServiceApi.md#RunWorkflow) | **Post** /runs | Run a workflow.



## CancelRun

> RunId CancelRun(ctx, runId)

Cancel a running workflow.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**runId** | **string**|  | 

### Return type

[**RunId**](RunId.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunLog

> RunLog GetRunLog(ctx, runId)

Get detailed info about a workflow run.

This endpoint provides detailed information about a given workflow run. The returned result has information about the outputs produced by this workflow (if available), a log object which allows the stderr and stdout to be retrieved, a log array so stderr/stdout for individual tasks can be retrieved, and the overall state of the workflow run (e.g. RUNNING, see the State section).

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**runId** | **string**|  | 

### Return type

[**RunLog**](RunLog.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunStatus

> RunStatus GetRunStatus(ctx, runId)

Get quick status info about a workflow run.

This provides an abbreviated (and likely fast depending on implementation) status of the running workflow, returning a simple result with the  overall state of the workflow run (e.g. RUNNING, see the State section).

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**runId** | **string**|  | 

### Return type

[**RunStatus**](RunStatus.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetServiceInfo

> ServiceInfo GetServiceInfo(ctx, )

Get information about Workflow Execution Service.

May include information related (but not limited to) the workflow descriptor formats, versions supported, the WES API versions supported, and information about general service availability.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListRuns

> RunListResponse ListRuns(ctx, optional)

List the workflow runs.

This list should be provided in a stable ordering. (The actual ordering is implementation dependent.) When paging through the list, the client should not make assumptions about live updates, but should assume the contents of the list reflect the workflow list at the moment that the first page is requested.  To monitor a specific workflow run, use GetRunStatus or GetRunLog.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListRunsOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ListRunsOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **optional.Int64**| OPTIONAL The preferred number of workflow runs to return in a page. If not provided, the implementation should use a default page size. The implementation must not return more items than &#x60;page_size&#x60;, but it may return fewer.  Clients should not assume that if fewer than &#x60;page_size&#x60; items are returned that all items have been returned.  The availability of additional pages is indicated by the value of &#x60;next_page_token&#x60; in the response. | 
 **pageToken** | **optional.String**| OPTIONAL Token to use to indicate where to start getting results. If unspecified, return the first page of results. | 

### Return type

[**RunListResponse**](RunListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RunWorkflow

> RunId RunWorkflow(ctx, optional)

Run a workflow.

This endpoint creates a new workflow run and returns a `RunId` to monitor its progress.  The `workflow_attachment` array may be used to upload files that are required to execute the workflow, including the primary workflow, tools imported by the workflow, other files referenced by the workflow, or files which are part of the input.  The implementation should stage these files to a temporary directory and execute the workflow from there. These parts must have a Content-Disposition header with a \"filename\" provided for each part.  Filenames may include subdirectories, but must not include references to parent directories with '..' -- implementations should guard against maliciously constructed filenames.  The `workflow_url` is either an absolute URL to a workflow file that is accessible by the WES endpoint, or a relative URL corresponding to one of the files attached using `workflow_attachment`.  The `workflow_params` JSON object specifies input parameters, such as input files.  The exact format of the JSON object depends on the conventions of the workflow language being used.  Input files should either be absolute URLs, or relative URLs corresponding to files uploaded using `workflow_attachment`.  The WES endpoint must understand and be able to access URLs supplied in the input.  This is implementation specific.  The `workflow_type` is the type of workflow language and must be \"CWL\" or \"WDL\" currently (or another alternative  supported by this WES instance).  The `workflow_type_version` is the version of the workflow language submitted and must be one supported by this WES instance.  See the `RunRequest` documentation for details about other fields.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***RunWorkflowOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a RunWorkflowOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowParams** | **optional.String**|  | 
 **workflowType** | **optional.String**|  | 
 **workflowTypeVersion** | **optional.String**|  | 
 **tags** | **optional.String**|  | 
 **workflowEngineParameters** | **optional.String**|  | 
 **workflowUrl** | **optional.String**|  | 
 **workflowAttachment** | **optional.Interface of []*os.File**|  | 

### Return type

[**RunId**](RunId.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

