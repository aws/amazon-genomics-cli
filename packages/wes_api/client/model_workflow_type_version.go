/*
 * Workflow Execution Service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package wes_client

// WorkflowTypeVersion - Available workflow types supported by a given instance of the service.
type WorkflowTypeVersion struct {

	// an array of one or more acceptable types for the `workflow_type`
	WorkflowTypeVersion []string `json:"workflow_type_version,omitempty"`
}
