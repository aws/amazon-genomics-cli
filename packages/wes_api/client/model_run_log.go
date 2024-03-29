/*
 * Workflow Execution Service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package wes_client

type RunLog struct {

	// workflow run ID
	RunId string `json:"run_id,omitempty"`

	Request RunRequest `json:"request,omitempty"`

	State State `json:"state,omitempty"`

	RunLog Log `json:"run_log,omitempty"`

	// The logs, and other key info like timing and exit code, for each step in the workflow run.
	TaskLogs []Log `json:"task_logs,omitempty"`

	// The outputs from the workflow run.
	Outputs map[string]interface{} `json:"outputs,omitempty"`
}
