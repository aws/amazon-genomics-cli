package cfn

import (
	"context"
	"regexp"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type Stack struct {
	Id           string
	Name         string
	Status       types.StackStatus
	StatusReason string
}

func (c Client) ListStacks(regexNameFilter *regexp.Regexp, statusFilter []types.StackStatus) ([]Stack, error) {
	output, err := c.cfn.ListStacks(context.Background(), &cloudformation.ListStacksInput{
		StackStatusFilter: statusFilter,
	})
	if err != nil {
		return nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}

	var stacks []Stack
	for _, stackSummary := range output.StackSummaries {
		stackName := aws.ToString(stackSummary.StackName)
		if regexNameFilter == nil || regexNameFilter.MatchString(stackName) {
			stacks = append(stacks, Stack{
				Id:           aws.ToString(stackSummary.StackId),
				Name:         stackName,
				Status:       stackSummary.StackStatus,
				StatusReason: aws.ToString(stackSummary.StackStatusReason),
			})
		}
	}

	return stacks, nil
}
