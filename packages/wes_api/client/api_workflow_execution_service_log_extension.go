/*
 * Additional methods for actually retrieving data pointed to by URLs in WES responses. 
 */

package wes_client

import (
	_context "context"
	_ioutil "io/ioutil"
	_nethttp "net/http"
	_neturl "net/url"
	"fmt"
	"io"
	"strings"
)

// Linger please
var (
	_ _context.Context
)

/*
GetRunLogData Get data linked to by GetRunLog.
Returns a stream for the content of a URL referenced in a GetRunLog response.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param runId The run ID used for the GetRunLog call
 * @param dataUrl The URL string from the GetRunLog call
@return *io.ReadCloser
*/
func (a *WorkflowExecutionServiceApiService) GetRunLogData(ctx _context.Context, runId string, dataUrl string) (*io.ReadCloser, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarReturnValue  *io.ReadCloser = nil
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/runs/{run_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"run_id"+"}", _neturl.PathEscape(parameterToString(runId, "")), -1)

	// Evaluate dataUrl relative to localVarPath and replace localVarPath
	base, err := _neturl.Parse(localVarPath)
    if err != nil {
        return localVarReturnValue, nil, err
    }
	evaluated, err := base.Parse(dataUrl)
	if err != nil {
        return localVarReturnValue, nil, err
    }
	if (evaluated.Scheme != base.Scheme && !strings.HasPrefix(evaluated.Scheme, "http")) {
		// This doesn't look like something we can fetch
		return localVarReturnValue, nil, fmt.Errorf("WES cannot be used to retrieve %s", dataUrl)
	}
	localVarPath = evaluated.String()

	// Request headers will go in here.
	localVarHeaderParams := make(map[string]string)
	// We don't use any of these, but we need them to invoke prepareRequest.
	files := make(map[string][]byte)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// We don't need any accept type choosing logic; we can only accept plain text.
	localVarHeaderParams["Accept"] = "text/plain"

	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, files)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	// Make the request
	localVarHTTPResponse, err := a.client.callAPI(r)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}
	// Be ready to return the body stream
	localVarReturnValue = localVarHTTPResponse.Body

	if localVarHTTPResponse.StatusCode >= 300 {
		// Something has gone wrong sever-side (and this isn't a redirect)
		// Fetch the entire body
		localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
		localVarHTTPResponse.Body.Close()
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if (err) {
			// Something went wrong during error download.
			// Add that error to our error.
			newErr.error = fmt.Errorf("Failed to download body of HTTP error %d %s response: %v", localVarHTTPResponse.StatusCode, localVarHTTPResponse.Status, err)
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		// Otherwise, we downloaded something. Maybe we can parse it as a WES-style JSON error.
		var v ErrorResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			// Nope, it's not a WES-style error.
			// Don't explain why it's not parseable, just pass along what the server said.
			newErr.error = fmt.Errorf("Instead of log data, server sent a %d %s error with content: %s", localVarHTTPResponse.StatusCode, localVarHTTPResponse.Status, localVarBody)
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		// Otherwise, it is a WES-style error we can understand (even if not a normally acceptable WES error code)
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}
	// TODO: handle redirects?

	return localVarReturnValue, localVarHTTPResponse, nil
}

