package ddb

import (
	"context"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	exp "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type WorkflowSummary struct {
	Name string
}

func (c *Client) ListWorkflows(ctx context.Context, project, user string) ([]WorkflowSummary, error) {
	pk := exp.Value(renderPartitionKey(project, user))
	skPref := renderWorkflowPrefix()
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
	}
	p := dynamodb.NewQueryPaginator(c.svc, input)
	uniqueWorkflows := make(map[string]WorkflowSummary)
	for p.HasMorePages() {
		output, err := p.NextPage(ctx)
		if err != nil {
			return nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
		}
		var records []WorkflowInstance
		err = attributevalue.UnmarshalListOfMaps(output.Items, &records)
		if err != nil {
			return nil, err
		}
		for _, record := range records {
			uniqueWorkflows[record.WorkflowName] = WorkflowSummary{Name: record.WorkflowName}
		}
	}
	var workflowSummaries []WorkflowSummary
	for _, summary := range uniqueWorkflows {
		workflowSummaries = append(workflowSummaries, summary)
	}

	return workflowSummaries, nil
}

func renderWorkflowPrefix() string {
	return "WORKFLOW#"
}
