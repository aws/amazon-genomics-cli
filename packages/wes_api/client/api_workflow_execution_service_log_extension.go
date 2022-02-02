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
	"strings"
)

// Linger please
var (
	_ _context.Context
)

/*
GetRunLogData Get data linked to by GetRunLog.
Returns the content of a URL reverenced in a GetRunLog response, if that URL is also on the WES server.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param runId The run ID used for the GetRunLog call
 * @param dataUrl The URL string from the GetRunLog call
@return string
*/
func (a *WorkflowExecutionServiceApiService) GetRunLogData(ctx _context.Context, runId string, dataUrl string) (string, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarReturnValue  string
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

	localVarHeaderParams := make(map[string]string)
	files := make(map[string][]byte)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"text/plain"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, files)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(r)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	// TODO: Handle odd encodings, unacceptable UTF-8, etc. somehow.
	localVarReturnValue = string(localVarBody)

	return localVarReturnValue, localVarHTTPResponse, nil
}

