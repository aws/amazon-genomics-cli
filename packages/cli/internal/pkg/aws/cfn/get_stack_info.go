package cfn

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"
	"github.com/rs/zerolog/log"
)

type StackInfo struct {
	Id      string
	Status  types.StackStatus
	Outputs map[string]string
	Tags    map[string]string
}

func (c Client) GetStackInfo(stackName string) (StackInfo, error) {
	input := &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}
	output, err := c.cfn.DescribeStacks(context.Background(), input)
	if err != nil {
		if isNotFoundError(err) {
			log.Debug().Msgf("Failed get info for stack '%s': %v", stackName, err)
			return StackInfo{}, StackDoesNotExistError
		}
		return StackInfo{}, err
	}
	if output == nil || len(output.Stacks) != 1 {
		return StackInfo{}, fmt.Errorf("unable to find stack '%s'", stackName)
	}

	stack := output.Stacks[0]
	return StackInfo{
		Id:      aws.ToString(stack.StackId),
		Status:  stack.StackStatus,
		Outputs: outputsToMap(stack.Outputs),
		Tags:    tagsToMap(stack.Tags),
	}, nil
}

func isNotFoundError(err error) bool {
	var ae smithy.APIError
	if errors.As(err, &ae) {
		return strings.Contains(ae.ErrorMessage(), "does not exist")
	}
	return false
}

func tagsToMap(tags []types.Tag) map[string]string {
	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return tagMap
}

func outputsToMap(outputs []types.Output) map[string]string {
	outputMap := make(map[string]string, len(outputs))
	for _, output := range outputs {
		outputMap[aws.ToString(output.OutputKey)] = aws.ToString(output.OutputValue)
	}
	return outputMap
}
