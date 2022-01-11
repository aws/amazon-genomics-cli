package cfn

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

const (
	testStackName1         = "test-stack-name-1"
	testStackId1           = "test-stack-id-1"
	testStackName2         = "test-stack-name-2"
	testStackId2           = "test-stack-id-2"
	testStackStatusReason1 = "test-stack-status-reason-1"
	testStackStatusReason2 = "test-stack-status-reason-2"
)

var (
	testStackSummaries = []types.StackSummary{
		{
			StackName:         aws.String(testStackName1),
			StackId:           aws.String(testStackId1),
			StackStatus:       types.StackStatusCreateComplete,
			StackStatusReason: aws.String(testStackStatusReason1),
		},
		{
			StackName:         aws.String(testStackName2),
			StackId:           aws.String(testStackId2),
			StackStatus:       types.StackStatusCreateFailed,
			StackStatusReason: aws.String(testStackStatusReason2),
		},
	}
)

func (m *CfnMock) ListStacks(ctx context.Context, input *cloudformation.ListStacksInput, opts ...func(*cloudformation.Options)) (*cloudformation.ListStacksOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*cloudformation.ListStacksOutput), err
	}
	return nil, err
}

func TestClient_ListStacks(t *testing.T) {
	testCases := map[string]struct {
		regexFilter    *regexp.Regexp
		statusFilter   []types.StackStatus
		setupMocks     func() *CfnMock
		expectedStacks []Stack
		expectedErr    error
	}{
		"no filters": {
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("ListStacks", context.Background(), &cloudformation.ListStacksInput{}).
					Return(&cloudformation.ListStacksOutput{StackSummaries: testStackSummaries}, nil)
				return cfnMock
			},
			expectedStacks: []Stack{
				{
					Name:         testStackName1,
					Id:           testStackId1,
					Status:       types.StackStatusCreateComplete,
					StatusReason: testStackStatusReason1,
				},
				{
					Name:         testStackName2,
					Id:           testStackId2,
					Status:       types.StackStatusCreateFailed,
					StatusReason: testStackStatusReason2,
				},
			},
		},
		"status filter with empty list": {
			statusFilter: []types.StackStatus{types.StackStatusCreateInProgress},
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("ListStacks", context.Background(), &cloudformation.ListStacksInput{
					StackStatusFilter: []types.StackStatus{types.StackStatusCreateInProgress},
				}).Return(&cloudformation.ListStacksOutput{StackSummaries: []types.StackSummary{}}, nil)
				return cfnMock
			},
		},
		"name filter": {
			regexFilter: regexp.MustCompile(`.*2$`),
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("ListStacks", context.Background(), &cloudformation.ListStacksInput{}).
					Return(&cloudformation.ListStacksOutput{StackSummaries: testStackSummaries}, nil)
				return cfnMock
			},
			expectedStacks: []Stack{
				{
					Name:         testStackName2,
					Id:           testStackId2,
					Status:       types.StackStatusCreateFailed,
					StatusReason: testStackStatusReason2,
				},
			},
		},
		"list error": {
			setupMocks: func() *CfnMock {
				cfnMock := new(CfnMock)
				cfnMock.On("ListStacks", context.Background(), &cloudformation.ListStacksInput{}).
					Return(nil, fmt.Errorf("list error"))
				return cfnMock
			},
			expectedErr: fmt.Errorf("list error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cfnMock := tc.setupMocks()
			client := Client{
				cfn: cfnMock,
			}

			stacks, err := client.ListStacks(tc.regexFilter, tc.statusFilter)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedStacks, stacks)
			cfnMock.AssertExpectations(t)
		})
	}
}
