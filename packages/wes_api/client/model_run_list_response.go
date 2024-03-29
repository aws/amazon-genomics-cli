/*
 * Workflow Execution Service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package wes_client

// RunListResponse - The service will return a RunListResponse when receiving a successful RunListRequest.
type RunListResponse struct {

	// A list of workflow runs that the service has executed or is executing. The list is filtered to only include runs that the caller has permission to see.
	Runs []RunStatus `json:"runs,omitempty"`

	// A token which may be supplied as `page_token` in workflow run list request to get the next page of results.  An empty string indicates there are no more items to return.
	NextPageToken string `json:"next_page_token,omitempty"`
}
