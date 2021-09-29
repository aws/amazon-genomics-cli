package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	exp "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (c *Client) ListWorkflowInstances(ctx context.Context, project, user string, limit int) ([]WorkflowInstance, error) {
	pk := exp.Value(renderPartitionKey(project, user))
	skPref := renderWorkflowRunPrefix()
	keyCondition := exp.Key(pkAttrName).Equal(pk).And(exp.Key(lsi2SkAttrName).BeginsWith(skPref))
	expression, err := exp.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		KeyConditionExpression:    expression.KeyCondition(),
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		IndexName:                 aws.String(Lsi2Name),
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

func renderWorkflowRunPrefix() string {
	return "RUN#CREATED#"
}
