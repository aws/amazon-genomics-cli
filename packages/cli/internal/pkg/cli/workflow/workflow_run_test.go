package workflow

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	wesmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/wes"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testProjectName         = "TestProject1"
	testLocalWorkflowName   = "TestLocalWorkflowName1"
	testS3WorkflowName      = "TestS3WorkflowName2"
	testInvalidWorkflowName = "TestInvalidWorkflowName2"
	testContext1Name        = "TestContext1"
	testContext2Name        = "TestContext2"
	testDataFileName        = "data.txt"
	testDataFileLocalUrl    = "path/to/" + testDataFileName
	testDataFileS3Url       = "s3://path/to/" + testDataFileName
	testInputKey            = "Workflow.variable"
	testInputLocal          = `{"` + testInputKey + `":"` + testDataFileLocalUrl + `"}`
	testInputS3             = `{"` + testInputKey + `":"` + testDataFileS3Url + `"}`
	testInputLocalToS3      = `{"` + testInputKey + `":"s3://TestOutputBucket/project/TestProject1/userid/bender123/data/Workflow.variable/data.txt"}`
	testWorkflowTypeLang    = "TypeLanguage"
	testWorkflowTypeVer     = "TypeVersion"
	testOutputBucket        = "TestOutputBucket"
	testWorkflowKey         = "project/" + testProjectName + "/userid/" + testUserId + "/context/" + testContext1Name + "/workflow/" + testLocalWorkflowName + "/workflow.zip"
	testDataKey             = "project/" + testProjectName + "/userid/" + testUserId + "/data/" + testInputKey + "/" + testDataFileName
	testWorkflowLocalUrl    = "workflow/path/file.wdl"
	testWorkflowS3Url       = "s3://workflow/path/file.wdl"
	testWorkflowInvalidUrl  = ":NotURL:"
	testCompressedTmpPath   = "/tmp/123/workflow_1343535"
	testArgsFileName        = "args.txt"
	testArgumentsPath       = "workflow/path/" + testArgsFileName
	testWesUrl              = "https://TestWesUrl.com/prod"
	testContext1Stack       = "Agc-Context-TestProject1-" + testUserId + "-TestContext1"
	testContext2Stack       = "Agc-Context-TestProject1-" + testUserId + "-TestContext2"
	testRun1Id              = "TestRun1Id"
	testRun2Id              = "TestRun2Id"
	testStackId             = "TestStackId"
	testUserId              = "bender123"
	testTmpAttachmentPath   = "/tmp/attachment.file"
	testProjectFileDir      = "/tmp/path/to/project"
	testErrorPrefix         = "unable to run workflow: "
)

var testWorkflowType = spec.WorkflowType{Language: testWorkflowTypeLang, Version: testWorkflowTypeVer}

type WorkflowRunTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockProjectClient *storagemocks.MockProjectClient
	mockConfigClient  *storagemocks.MockConfigClient
	mockSsmClient     *awsmocks.MockSsmClient
	mockCfn           *awsmocks.MockCfnClient
	mockS3Client      *awsmocks.MockS3Client
	mockDdb           *awsmocks.MockDdbClient
	mockStorageClient *storagemocks.MockStorageClient
	mockOs            *iomocks.MockOS
	mockZip           *iomocks.MockZip
	mockTmp           *iomocks.MockTmp
	mockWes           *wesmocks.MockWesClient
	origRemoveFile    func(name string) error
	origCompressToTmp func(srcPath string) (string, error)
	origWriteToTmp    func(namePattern, content string) (string, error)

	testProjSpec         spec.Project
	absDataFilePath      string
	localWorkflowAbsPath string
	wfInstance           ddb.WorkflowInstance
	testStackInfo        cfn.StackInfo

	manager *Manager
}

func (s *WorkflowRunTestSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.mockProjectClient = storagemocks.NewMockProjectClient(s.ctrl)
	s.mockConfigClient = storagemocks.NewMockConfigClient(s.ctrl)
	s.mockSsmClient = awsmocks.NewMockSsmClient(s.ctrl)
	s.mockCfn = awsmocks.NewMockCfnClient(s.ctrl)
	s.mockS3Client = awsmocks.NewMockS3Client(s.ctrl)
	s.mockDdb = awsmocks.NewMockDdbClient(s.ctrl)
	s.mockStorageClient = storagemocks.NewMockStorageClient(s.ctrl)
	s.mockOs = iomocks.NewMockOS(s.ctrl)
	s.mockZip = iomocks.NewMockZip(s.ctrl)
	s.mockTmp = iomocks.NewMockTmp(s.ctrl)
	s.mockWes = wesmocks.NewMockWesClient(s.ctrl)

	s.origRemoveFile, removeFile = removeFile, s.mockOs.Remove
	s.origCompressToTmp, compressToTmp = compressToTmp, s.mockZip.CompressToTmp
	s.origWriteToTmp, writeToTmp = writeToTmp, s.mockTmp.Write

	// data file path is relative to inputs file (usually the workflow folder)
	absDataFilePath, err := filepath.Abs(filepath.Join(filepath.Dir(testArgumentsPath), testDataFileLocalUrl))
	require.NoError(s.T(), err)
	s.absDataFilePath = absDataFilePath
	s.localWorkflowAbsPath = filepath.Join(testProjectFileDir, testWorkflowLocalUrl)

	s.manager = &Manager{
		Project:    s.mockProjectClient,
		Ssm:        s.mockSsmClient,
		Cfn:        s.mockCfn,
		S3:         s.mockS3Client,
		Ddb:        s.mockDdb,
		Storage:    s.mockStorageClient,
		Config:     s.mockConfigClient,
		WesFactory: func(_ string) (wes.Interface, error) { return s.mockWes, nil },
	}

	s.testProjSpec = spec.Project{
		Name: testProjectName,
		Workflows: map[string]spec.Workflow{
			testLocalWorkflowName: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowLocalUrl,
			},
			testS3WorkflowName: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowS3Url,
			},
			testInvalidWorkflowName: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowInvalidUrl,
			},
		},
		Contexts: map[string]spec.Context{
			testContext1Name: {
				Engines: []spec.Engine{
					{
						Type:   "wdl",
						Engine: "cromwell",
					},
				},
			},
		},
	}

	s.wfInstance = ddb.WorkflowInstance{
		RunId:        testRun1Id,
		WorkflowName: testLocalWorkflowName,
		ContextName:  testContext1Name,
		ProjectName:  testProjectName,
		UserId:       testUserId,
	}

	s.testStackInfo = cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWesUrl},
	}
}

func (s *WorkflowRunTestSuite) AfterTest(_, _ string) {
	removeFile = s.origRemoveFile
	compressToTmp = s.origCompressToTmp
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_LocalFile_WithS3Args() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	uploadCall := s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputS3), nil)
	s.mockTmp.EXPECT().Write(testArgsFileName+"_*", testInputS3).Return(testTmpAttachmentPath, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).After(uploadCall).Return(nil)
	s.mockOs.EXPECT().Remove(testTmpAttachmentPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath)
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_LocalFile_WithLocalArgs() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	uploadCall := s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputLocal), nil)
	s.mockTmp.EXPECT().Write(testArgsFileName+"_*", testInputLocalToS3).Return(testTmpAttachmentPath, nil)
	s.mockS3Client.EXPECT().SyncFile(testOutputBucket, testDataKey, s.absDataFilePath).Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).After(uploadCall).Return(nil)
	s.mockOs.EXPECT().Remove(testTmpAttachmentPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath)
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_LocalFile_NoArgs() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	uploadCall := s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).After(uploadCall).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_S3Object_WithLocalArgs() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputLocal), nil)
	s.mockTmp.EXPECT().Write(testArgsFileName+"_*", testInputLocalToS3).Return(testTmpAttachmentPath, nil)
	s.mockS3Client.EXPECT().SyncFile(testOutputBucket, testDataKey, s.absDataFilePath).Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.wfInstance.WorkflowName = testS3WorkflowName
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testTmpAttachmentPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testS3WorkflowName, testArgumentsPath)
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_S3Object_NoArgs() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.wfInstance.WorkflowName = testS3WorkflowName
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	actualId, err := s.manager.RunWorkflow(testContext1Name, testS3WorkflowName, "")
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_ReadProjectSpecFailure() {
	errorMessage := "failed to read project specification"
	s.mockProjectClient.EXPECT().Read().Return(spec.Project{}, errors.New(errorMessage))

	actualId, err := s.manager.RunWorkflow(testContext1Name, "", "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_MissingWorkflowSpec() {
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, "dummy", "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+"workflow 'dummy' is not defined in Project 'TestProject1' specification")
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_InvalidWorkflowDefinitionUrl() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testInvalidWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+`parse ":NotURL:": missing protocol scheme`)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_CompressionFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot compress file"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return("", errors.New(errorMessage))

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_SSMClientFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot connect to SSM"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return("", errors.New(errorMessage))

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_UploadToS3Failed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot upload to S3"
	expectedInfix := "unable to upload s3://TestOutputBucket/project/TestProject1/userid/bender123/context/TestContext1/workflow/TestLocalWorkflowName1/workflow.zip: "
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+expectedInfix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_ReadArgsFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot read input"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte{}, errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_UploadInputFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot upload input"
	expectedInfix := "unable to sync s3://TestOutputBucket/project/TestProject1/userid/bender123/data/Workflow.variable/data.txt: "
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputLocal), nil)
	s.mockS3Client.EXPECT().SyncFile(testOutputBucket, testDataKey, s.absDataFilePath).Return(errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+expectedInfix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_CfnFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot call CFN"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{}, errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_CfnMissingWesUrlFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "wes endpoint for workflow type 'TypeLanguage' is missing in engine stack 'TestStackId'"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	stackInfo := cfn.StackInfo{
		Id:      testStackId,
		Outputs: map[string]string{},
	}
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(stackInfo, nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_WesFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot call WES"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockZip.EXPECT().CompressToTmp(s.localWorkflowAbsPath).Return(testCompressedTmpPath, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowKey, testCompressedTmpPath).Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return("", errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_DeployValidationFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	errorMessage := "context 'TestContext1' is not deployed"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, cfn.StackDoesNotExistError)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_DeployValidationCfnFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	errorMessage := "some cfn error"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, errors.New(errorMessage))

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func TestWorkflowRunTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowRunTestSuite))
}
