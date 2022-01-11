package ddb

import (
	"context"
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	exp "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (c *Client) ListWorkflowInstancesByName(ctx context.Context, project, user, workflowName string, limit int) ([]WorkflowInstance, error) {
	pk := exp.Value(renderPartitionKey(project, user))
	skPref := renderWorkflowNamePrefix(workflowName)
	keyCondition := exp.Key(pkAttrName).Equal(pk).And(exp.Key(lsi1SkAttrName).BeginsWith(skPref))
	expression, err := exp.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		KeyConditionExpression:    expression.KeyCondition(),
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		IndexName:                 aws.String(Lsi1Name),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     aws.Int32(int32(limit)),
	}
	output, err := c.svc.Query(ctx, input)
	if err != nil {
		return nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	var instances []WorkflowInstance
	err = attributevalue.UnmarshalListOfMaps(output.Items, &instances)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func renderWorkflowNamePrefix(workflowName string) string {
	return fmt.Sprintf("RUN#WORKFLOW#%s#", workflowName)
}
