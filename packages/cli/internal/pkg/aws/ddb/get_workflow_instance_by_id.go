package ddb

import (
	"context"
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (c *Client) GetWorkflowInstanceById(ctx context.Context, project, user, runId string) (WorkflowInstance, error) {
	pk := renderPartitionKey(project, user)
	sk := renderRunSortKey(runId)
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			pkAttrName: &types.AttributeValueMemberS{Value: pk},
			skAttrName: &types.AttributeValueMemberS{Value: sk},
		},
		TableName: aws.String(TableName),
	}
	output, err := c.svc.GetItem(ctx, input)
	if err != nil {
		return WorkflowInstance{}, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	if output.Item == nil {
		return WorkflowInstance{}, fmt.Errorf("workflow instance with id '%s' does not exist", runId)
	}
	var instance WorkflowInstance
	if err := attributevalue.UnmarshalMap(output.Item, &instance); err != nil {
		return WorkflowInstance{}, err
	}
	return instance, nil
}
