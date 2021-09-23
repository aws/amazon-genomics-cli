package ddb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var now = time.Now

type WorkflowInstance struct {
	RunId        string
	WorkflowName string
	ContextName  string
	ProjectName  string
	UserId       string
	CreatedTime  string
}

func (c *Client) WriteWorkflowInstance(ctx context.Context, instance WorkflowInstance) error {
	instance.CreatedTime = now().UTC().Format(time.RFC3339)
	record, err := createInstanceRecord(instance)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      record,
		TableName: aws.String(TableName),
	}
	if _, err := c.svc.PutItem(ctx, input); err != nil {
		return err
	}
	return nil
}

func createInstanceRecord(instance WorkflowInstance) (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(instance)
	if err != nil {
		return nil, err
	}
	av[pkAttrName] = &types.AttributeValueMemberS{
		Value: renderPartitionKey(instance.ProjectName, instance.UserId),
	}
	av[skAttrName] = &types.AttributeValueMemberS{Value: renderRunSortKey(instance.RunId)}
	av[lsi1SkAttrName] = &types.AttributeValueMemberS{Value: renderRunLsi1SortKey(instance)}
	av[lsi2SkAttrName] = &types.AttributeValueMemberS{Value: renderRunLsi2SortKey(instance)}
	av[lsi3SkAttrName] = &types.AttributeValueMemberS{Value: renderRunLsi3SortKey(instance)}
	return av, nil
}

func renderPartitionKey(project, user string) string {
	return fmt.Sprintf("PROJECT#%s#USER#%s", project, user)
}

func renderRunSortKey(runId string) string {
	return fmt.Sprintf("RUN#%s", runId)
}

func renderRunLsi1SortKey(instance WorkflowInstance) string {
	return fmt.Sprintf("RUN#WORKFLOW#%s#CREATED#%s", instance.WorkflowName, instance.CreatedTime)
}

func renderRunLsi2SortKey(instance WorkflowInstance) string {
	return fmt.Sprintf("RUN#CREATED#%s", instance.CreatedTime)
}

func renderRunLsi3SortKey(instance WorkflowInstance) string {
	return fmt.Sprintf("RUN#CONTEXT#%s#CREATED#%s", instance.ContextName, instance.CreatedTime)
}
