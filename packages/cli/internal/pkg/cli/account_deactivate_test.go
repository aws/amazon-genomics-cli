package cli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/stretchr/testify/assert"
)

const (
	testDeactivateStackName1 = "test-deactivate-stack-name-1"
	testDeactivateStackId1   = "test-deactivate-stack-id-1"
	testDeactivateStackName2 = "test-deactivate-stack-name-2"
	testDeactivateStackId2   = "test-deactivate-stack-id-2"
	testDeactivateStackName3 = "test-deactivate-stack-name-3"
	testDeactivateStackId3   = "test-deactivate-stack-id-3"
)

var (
	testDeactivateStack1 = cfn.Stack{
		Name: testDeactivateStackName1,
		Id:   testDeactivateStackId1,
	}
	testDeactivateStack2 = cfn.Stack{
		Name: testDeactivateStackName2,
		Id:   testDeactivateStackId2,
	}
	testDeactivateStack3 = cfn.Stack{
		Name: testDeactivateStackName3,
		Id:   testDeactivateStackId3,
	}
)

func TestAccountDeactivateOpts_Load(t *testing.T) {
	testCases := map[string]struct {
		expectedStacks []cfn.Stack
		setupMocks     func(*testing.T) mockClients
		expectedErr    error
	}{
		"success": {
			expectedStacks: []cfn.Stack{testDeactivateStack1},
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.cfnMock.EXPECT().ListStacks(regexp.MustCompile(`^Agc-.*$`), cfn.ActiveStacksFilter).
					Return([]cfn.Stack{testDeactivateStack1, testDeactivateStack2}, nil)
				mocks.cfnMock.EXPECT().GetStackTags(testDeactivateStackId1).Return(map[string]string{"application-name": "agc"}, nil)
				mocks.cfnMock.EXPECT().GetStackTags(testDeactivateStackId2).Return(map[string]string{}, nil)
				return mocks
			},
		},
		"list error": {
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.cfnMock.EXPECT().ListStacks(regexp.MustCompile(`^Agc-.*$`), cfn.ActiveStacksFilter).
					Return(nil, fmt.Errorf("some list error"))
				return mocks
			},
			expectedErr: fmt.Errorf("An error occurred while deactivating the account. Error was: 'some list error'"),
		},
		"get tags error": {
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.cfnMock.EXPECT().ListStacks(regexp.MustCompile(`^Agc-.*$`), cfn.ActiveStacksFilter).
					Return([]cfn.Stack{testDeactivateStack1, testDeactivateStack2}, nil)
				mocks.cfnMock.EXPECT().GetStackTags(testDeactivateStackId1).Return(nil, fmt.Errorf("some tags error"))
				return mocks
			},
			expectedErr: fmt.Errorf("An error occurred while deactivating the account. Error was: 'some tags error'"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mocks := tc.setupMocks(t)
			defer mocks.ctrl.Finish()
			opts := &accountDeactivateOpts{
				cfnClient: mocks.cfnMock,
			}

			err := opts.LoadStacks()
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedStacks, opts.stacks)
		})
	}
}
func TestAccountDeactivateOpts_Validate(t *testing.T) {
	testCases := map[string]struct {
		force       bool
		stacks      []cfn.Stack
		expectedErr bool
	}{
		"single stack without force no error": {
			stacks: []cfn.Stack{testDeactivateStack1},
		},
		"two stacks without force no error": {
			stacks: []cfn.Stack{testDeactivateStack1, testDeactivateStack2},
		},
		"many stacks with force no error": {
			force:  true,
			stacks: []cfn.Stack{testDeactivateStack1, testDeactivateStack2, testDeactivateStack3},
		},
		"too many stacks without force error": {
			stacks:      []cfn.Stack{testDeactivateStack1, testDeactivateStack2, testDeactivateStack3},
			expectedErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			opts := &accountDeactivateOpts{
				accountDeactivateVars: accountDeactivateVars{
					force: tc.force,
				},
				stacks: tc.stacks,
			}

			err := opts.Validate()
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestAccountDeactivateOpts_Execute(t *testing.T) {
	testCases := map[string]struct {
		setupMocks  func(*testing.T) mockClients
		expectedErr error
	}{
		"delete success": {
			setupMocks: func(t *testing.T) mockClients {
				testChan1 := make(chan cfn.DeletionResult)
				testChan2 := make(chan cfn.DeletionResult)
				mocks := createMocks(t)
				mocks.cfnMock.EXPECT().DeleteStack(testDeactivateStackId1).Return(testChan1, nil)
				mocks.cfnMock.EXPECT().DeleteStack(testDeactivateStackId2).Return(testChan1, nil)
				go func() { testChan1 <- cfn.DeletionResult{}; close(testChan1) }()
				go func() { testChan2 <- cfn.DeletionResult{}; close(testChan2) }()
				return mocks
			},
		},
		"delete failure": {
			setupMocks: func(t *testing.T) mockClients {
				testChan1 := make(chan cfn.DeletionResult)
				testChan2 := make(chan cfn.DeletionResult)
				mocks := createMocks(t)
				mocks.cfnMock.EXPECT().DeleteStack(testDeactivateStackId1).Return(testChan1, nil)
				mocks.cfnMock.EXPECT().DeleteStack(testDeactivateStackId2).Return(testChan2, nil)
				go func() { testChan1 <- cfn.DeletionResult{}; close(testChan1) }()
				go func() { testChan2 <- cfn.DeletionResult{Error: fmt.Errorf("some delete error")}; close(testChan2) }()
				return mocks
			},
			expectedErr: fmt.Errorf("failed to delete stack '%s: %w", testDeactivateStack2.Name, fmt.Errorf("some delete error")),
		},
		"delete error": {
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.cfnMock.EXPECT().DeleteStack(testDeactivateStackId1).Return(nil, fmt.Errorf("some delete error"))
				return mocks
			},
			expectedErr: fmt.Errorf("some delete error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mocks := tc.setupMocks(t)
			defer mocks.ctrl.Finish()
			opts := &accountDeactivateOpts{
				stacks:    []cfn.Stack{testDeactivateStack1, testDeactivateStack2},
				cfnClient: mocks.cfnMock,
			}

			err := opts.Execute()
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
