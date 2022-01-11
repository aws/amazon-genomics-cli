package cfn

import (
	"context"
	"fmt"
	"testing"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

func (m *CfnMock) DeleteStack(ctx context.Context, input *cloudformation.DeleteStackInput, opts ...func(*cloudformation.Options)) (*cloudformation.DeleteStackOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*cloudformation.DeleteStackOutput), err
	}
	return nil, err
}

func TestClient_DeleteStack(t *testing.T) {
	realSleepDuration := sleepDuration
	sleepDuration = time.Millisecond
	defer func() { sleepDuration = realSleepDuration }()

	testCases := map[string]struct {
		setupMocks      func() *CfnMock
		expectedTracker DeletionResult
		expectedErr     error
	}{
		"delete success": {
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("DeleteStack", context.Background(), &cloudformation.DeleteStackInput{
					StackName: aws.String(testStackName),
				}).Return(nil, nil)
				cfnMock.On("DescribeStacks", context.Background(), &cloudformation.DescribeStacksInput{
					StackName: aws.String(testStackName),
				}).Return(&cloudformation.DescribeStacksOutput{Stacks: []types.Stack{{StackStatus: types.StackStatusDeleteInProgress}}}, nil).Once()
				cfnMock.On("DescribeStacks", context.Background(), &cloudformation.DescribeStacksInput{
					StackName: aws.String(testStackName),
				}).Return(&cloudformation.DescribeStacksOutput{Stacks: []types.Stack{{StackStatus: types.StackStatusDeleteComplete}}}, nil)
				return cfnMock
			},
			expectedTracker: DeletionResult{},
		},
		"delete failure": {
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("DeleteStack", context.Background(), &cloudformation.DeleteStackInput{
					StackName: aws.String(testStackName),
				}).Return(nil, nil)
				cfnMock.On("DescribeStacks", context.Background(), &cloudformation.DescribeStacksInput{
					StackName: aws.String(testStackName),
				}).Return(&cloudformation.DescribeStacksOutput{Stacks: []types.Stack{{StackStatus: types.StackStatusDeleteFailed}}}, nil)
				return cfnMock
			},
			expectedTracker: DeletionResult{Error: fmt.Errorf("failed to delete stack")},
		},
		"unknown status": {
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("DeleteStack", context.Background(), &cloudformation.DeleteStackInput{
					StackName: aws.String(testStackName),
				}).Return(nil, nil)
				cfnMock.On("DescribeStacks", context.Background(), &cloudformation.DescribeStacksInput{
					StackName: aws.String(testStackName),
				}).Return(&cloudformation.DescribeStacksOutput{Stacks: []types.Stack{{StackStatus: types.StackStatusUpdateComplete}}}, nil)
				return cfnMock
			},
			expectedTracker: DeletionResult{Error: fmt.Errorf("unexpected status of the stack: %s", types.StackStatusUpdateComplete)},
		},
		"delete error": {
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("DeleteStack", context.Background(), &cloudformation.DeleteStackInput{
					StackName: aws.String(testStackName),
				}).Return(nil, fmt.Errorf("some delete error"))
				return cfnMock
			},
			expectedErr: fmt.Errorf("some delete error"),
		},
		"status error": {
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("DeleteStack", context.Background(), &cloudformation.DeleteStackInput{
					StackName: aws.String(testStackName),
				}).Return(nil, nil)
				cfnMock.On("DescribeStacks", context.Background(), &cloudformation.DescribeStacksInput{
					StackName: aws.String(testStackName),
				}).Return(nil, fmt.Errorf("some status error"))
				return cfnMock
			},
			expectedTracker: DeletionResult{Error: fmt.Errorf("some status error")},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cfnMock := tc.setupMocks()
			client := Client{
				cfn: cfnMock,
			}

			tracker, err := client.DeleteStack(testStackName)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTracker, <-tracker)
			}
			cfnMock.AssertExpectations(t)
		})
	}
}
