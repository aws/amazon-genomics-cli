package cfn

import (
	"context"
	"fmt"
	"testing"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

const (
	testStackName   = "test-stack-name"
	testStackId     = "test-stack-id"
	testOutputKey   = "test-output-key"
	testOutputValue = "test-output-value"
	testTagKey      = "test-tag-key"
	testTagValue    = "test-tag-value"
)

var (
	testDescribeInput = &cloudformation.DescribeStacksInput{StackName: aws.String(testStackName)}
)

func (m *CfnMock) DescribeStacks(ctx context.Context, input *cloudformation.DescribeStacksInput, opts ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*cloudformation.DescribeStacksOutput), err
	}
	return nil, err
}

func TestClient_GetStackInfo(t *testing.T) {
	client := NewMockClient()
	describeStacksOutput := &cloudformation.DescribeStacksOutput{
		Stacks: []types.Stack{{
			StackId:     aws.String(testStackId),
			StackStatus: types.StackStatusUpdateComplete,
			Outputs:     []types.Output{{OutputKey: aws.String(testOutputKey), OutputValue: aws.String(testOutputValue)}},
			Tags:        []types.Tag{{Key: aws.String(testTagKey), Value: aws.String(testTagValue)}}}},
	}
	client.cfn.(*CfnMock).On("DescribeStacks", context.Background(), testDescribeInput).
		Return(describeStacksOutput, nil)
	info, err := client.GetStackInfo(testStackName)
	assert.NoError(t, err)
	assert.Equal(t, StackInfo{
		Id:      testStackId,
		Status:  types.StackStatusUpdateComplete,
		Outputs: map[string]string{testOutputKey: testOutputValue},
		Tags:    map[string]string{testTagKey: testTagValue},
	}, info)
	client.cfn.(*CfnMock).AssertExpectations(t)
}

func TestClient_GetStackInfo_WithDescribeError(t *testing.T) {
	client := NewMockClient()
	describeStacksError := fmt.Errorf("some describe error")
	client.cfn.(*CfnMock).On("DescribeStacks", context.Background(), testDescribeInput).
		Return(nil, describeStacksError)
	_, err := client.GetStackInfo(testStackName)
	assert.Error(t, err, describeStacksError)
	client.cfn.(*CfnMock).AssertExpectations(t)
}

func TestClient_GetStackInfo_WithNoStackError(t *testing.T) {
	client := NewMockClient()
	describeStacksOutput := &cloudformation.DescribeStacksOutput{}
	client.cfn.(*CfnMock).On("DescribeStacks", context.Background(), testDescribeInput).
		Return(describeStacksOutput, nil)
	_, err := client.GetStackInfo(testStackName)
	assert.Error(t, err, fmt.Errorf("unable to find stack '%s'", testStackName))
	client.cfn.(*CfnMock).AssertExpectations(t)
}
