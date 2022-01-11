package cfn

import (
	"context"
	"regexp"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type Interface interface {
	GetStackInfo(stackName string) (StackInfo, error)
	GetStackOutputs(stackName string) (map[string]string, error)
	GetStackStatus(stackName string) (types.StackStatus, error)
	GetStackTags(stackName string) (map[string]string, error)
	ListStacks(regexNameFilter *regexp.Regexp, statusFilter []types.StackStatus) ([]Stack, error)
	DeleteStack(stackId string) (chan DeletionResult, error)
}

type cfnInterface interface {
	cloudformation.DescribeStacksAPIClient
	cloudformation.ListStacksAPIClient
	DeleteStack(ctx context.Context, input *cloudformation.DeleteStackInput, opts ...func(*cloudformation.Options)) (*cloudformation.DeleteStackOutput, error)
}
