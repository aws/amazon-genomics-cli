package context

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/environment"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestManager_Deploy(t *testing.T) {
	origVerbose := logging.Verbose
	origDisplayProgressBar := displayProgressBar
	defer func() {
		logging.Verbose = origVerbose
		displayProgressBar = origDisplayProgressBar
	}()
	logging.Verbose = false
	contextList := []string{testContextName1}

	testCases := map[string]struct {
		setupMocks                 func(*testing.T) mockClients
		expectedProgressResultList []ProgressResult
		contextList                []string
	}{
		"deploy success": {
			contextList: []string{testContextName3},
			expectedProgressResultList: []ProgressResult{
				{Outputs: []string{"some message"}, Context: testContextName3},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				defer close(mockClients.progressStream1)
				defer close(mockClients.progressStream2)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				mockClients.ecrClientMock.EXPECT().VerifyImageExists(environment.CommonImages["NEXTFLOW"]).Return(nil)
				clearContext := mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				mockClients.cdkMock.EXPECT().DeployApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Len(39), testContextName3).After(clearContext).Return(mockClients.progressStream1, nil)
				displayProgressBar = mockClients.cdkMock.DisplayProgressBar
				mockClients.cdkMock.EXPECT().DisplayProgressBar(fmt.Sprintf("Deploying resources for context(s) %s", []string{testContextName3}), []cdk.ProgressStream{mockClients.progressStream1}).Return([]cdk.Result{{Outputs: []string{"some message"}, ExecutionName: testContextName3}})
				return mockClients
			},
		},
		"multiple deploy success": {
			contextList: []string{testContextName1, testContextName2},
			expectedProgressResultList: []ProgressResult{
				{Outputs: []string{"some message"}, Context: testContextName1},
				{Outputs: []string{"some other message"}, Context: testContextName2},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				defer close(mockClients.progressStream1)
				defer close(mockClients.progressStream2)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Times(2).Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Times(2).Return(testArtifactBucket, nil)
				mockClients.ecrClientMock.EXPECT().VerifyImageExists(environment.CommonImages["CROMWELL"]).Times(2).Return(nil)
				clearContext := mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				clearContext2 := mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				mockClients.cdkMock.EXPECT().DeployApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Len(39), testContextName1).After(clearContext).Return(mockClients.progressStream1, nil)
				mockClients.cdkMock.EXPECT().DeployApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Len(39), testContextName2).After(clearContext2).Return(mockClients.progressStream2, nil)
				displayProgressBar = mockClients.cdkMock.DisplayProgressBar
				expectedCdkResult := []cdk.Result{{Outputs: []string{"some message"}, ExecutionName: testContextName1}, {Outputs: []string{"some other message"}, ExecutionName: testContextName2}}
				mockClients.cdkMock.EXPECT().DisplayProgressBar(fmt.Sprintf("Deploying resources for context(s) %s", []string{testContextName1, testContextName2}), []cdk.ProgressStream{mockClients.progressStream1, mockClients.progressStream2}).Return(expectedCdkResult)
				return mockClients
			},
		},
		"image does not exist": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("some error occurred"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				defer close(mockClients.progressStream1)
				defer close(mockClients.progressStream2)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				mockClients.ecrClientMock.EXPECT().VerifyImageExists(environment.CommonImages["CROMWELL"]).Return(fmt.Errorf("some error occurred"))
				return mockClients
			},
		},
		"engine name does not exist": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: actionableerror.New(
					fmt.Errorf("the engine name in your context file 'engine' does not exist"),
					"Please check your agc config file for the engine you have supplied",
				), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				defer close(mockClients.progressStream1)
				defer close(mockClients.progressStream2)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)

				testInvalidEngineProjectSpec := spec.Project{
					Name: testProjectName,
					Data: []spec.Data{{Location: testS3Location1}, {Location: testS3Location2, ReadOnly: true}},
					Contexts: map[string]spec.Context{
						testContextName1: {
							Engines: []spec.Engine{
								{Type: "wdl", Engine: "badEngineName"},
							},
						},
						testContextName2: {
							Engines: []spec.Engine{
								{Type: "wdl", Engine: "badEngineName"},
							},
						},
						testContextName3: {
							Engines: []spec.Engine{
								{Type: "wdl", Engine: "badEngineName"},
							},
						},
					},
				}
				mockClients.projMock.EXPECT().Read().Return(testInvalidEngineProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				return mockClients
			},
		},
		"read error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("some read error"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{}, fmt.Errorf("some read error"))
				return mockClients
			},
		},
		"output bucket error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("some outbut bucket error"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return("", fmt.Errorf("some outbut bucket error"))
				return mockClients
			},
		},
		"context error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("context 'testContextName1' is not defined in Project 'testProjectName' specification"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				projSpec := testValidProjectSpec
				projSpec.Contexts = nil
				mockClients.projMock.EXPECT().Read().Return(projSpec, nil)
				return mockClients
			},
		},
		"artifact bucket error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("some artifact bucket error"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return("", fmt.Errorf("some artifact bucket error"))
				return mockClients
			},
		},
		"deploy error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("some context error"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				mockClients.cdkMock.EXPECT().DeployApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Len(39), testContextName1).Return(nil, fmt.Errorf("some context error"))
				mockClients.ecrClientMock.EXPECT().VerifyImageExists(environment.CommonImages["CROMWELL"]).Return(nil)
				return mockClients
			},
		},
		"clear context error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("failed to clear context"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(fmt.Errorf("failed to clear context"))
				return mockClients
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockClients := tc.setupMocks(t)
			defer mockClients.ctrl.Finish()
			manager := Manager{
				Cdk:       mockClients.cdkMock,
				Project:   mockClients.projMock,
				Ssm:       mockClients.ssmMock,
				Config:    mockClients.configMock,
				Cfn:       mockClients.cfnMock,
				ecrClient: mockClients.ecrClientMock,
				baseProps: baseProps{homeDir: testHomeDir},
				imageRefs: environment.CommonImages,
				region:    "us-east-1",
			}

			progressResultList := manager.Deploy(tc.contextList)

			if len(progressResultList) != len(tc.expectedProgressResultList) {
				assert.Equal(t, tc.expectedProgressResultList, progressResultList)
			}

			for i, progressResult := range progressResultList {
				expectedProgressResult := tc.expectedProgressResultList[i]
				assert.Equal(t, expectedProgressResult.Context, progressResult.Context)
				if expectedProgressResult.Err != nil {
					assert.Error(t, progressResult.Err, expectedProgressResult)
				} else {
					assert.NoError(t, progressResult.Err)
				}
			}
		})
	}
}
