package ddb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Interface interface {
	WriteWorkflowInstance(ctx context.Context, instance WorkflowInstance) error
	ListWorkflows(ctx context.Context, project, user string) ([]WorkflowSummary, error)
	ListWorkflowInstancesByName(ctx context.Context, project, user, workflowName string, limit int) ([]WorkflowInstance, error)
	ListWorkflowInstancesByContext(ctx context.Context, project, user, contextName string, limit int) ([]WorkflowInstance, error)
	ListWorkflowInstances(ctx context.Context, project, user string, limit int) ([]WorkflowInstance, error)
	GetWorkflowInstanceById(ctx context.Context, project, user, runId string) (WorkflowInstance, error)
}

type ApiInterface interface {
	dynamodb.QueryAPIClient
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}
