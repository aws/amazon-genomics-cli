package ddb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	exp "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (c *Client) ListWorkflowInstancesByContext(ctx context.Context, project, user, contextName string, limit int) ([]WorkflowInstance, error) {
	pk := exp.Value(renderPartitionKey(project, user))
	skPref := renderContextNamePrefix(contextName)
	keyCondition := exp.Key(pkAttrName).Equal(pk).And(exp.Key(lsi3SkAttrName).BeginsWith(skPref))
	expression, err := exp.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		KeyConditionExpression:    expression.KeyCondition(),
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		IndexName:                 aws.String(Lsi3Name),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     aws.Int32(int32(limit)),
	}
	output, err := c.svc.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	var instances []WorkflowInstance
	err = attributevalue.UnmarshalListOfMaps(output.Items, &instances)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func renderContextNamePrefix(contextName string) string {
	return fmt.Sprintf("RUN#CONTEXT#%s#", contextName)
}
