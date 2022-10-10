package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	testProjectName          = "TestProject1"
	testLocalWorkflowName    = "TestLocalWorkflowName1"
	testS3WorkflowName       = "TestS3WorkflowName2"
	testInvalidWorkflowName  = "TestInvalidWorkflowName2"
	testContext1Name         = "TestContext1"
	testContext2Name         = "TestContext2"
	testDataFileName         = "data.txt"
	testDataFileLocalUrl     = "path/to/" + testDataFileName
	testDataFileS3Url        = "s3://path/to/" + testDataFileName
	testInputKey             = "Workflow.variable"
	testInputLocal           = `{"` + testInputKey + `":"` + testDataFileLocalUrl + `"}`
	testInputS3              = `{"` + testInputKey + `":"` + testDataFileS3Url + `"}`
	testInputLocalToS3       = `{"` + testInputKey + `":"s3://TestOutputBucket/project/TestProject1/userid/bender123/data/Workflow.variable/data.txt"}`
	testOptionFileLocal      = `{"testOptionName": "testOption"}`
	testWorkflowTypeLang     = "TypeLanguage"
	testWorkflowTypeVer      = "TypeVersion"
	testOutputBucket         = "TestOutputBucket"
	testWorkflowKey          = "project/" + testProjectName + "/userid/" + testUserId + "/context/" + testContext1Name + "/workflow/" + testLocalWorkflowName
	testWorkflowZipKey       = testWorkflowKey + "/workflow.zip"
	testWorkflowLocalUrl     = "workflow/path/file.wdl"
	testMANIFESTPath         = "extra/MANIFEST.json"
	testMANIFEST             = `{"mainWorkflowURL": "haplotypecaller-gvcf-gatk4.wdl","inputFileURLs": ["haplotypecaller-gvcf-gatk4.hg38.wgs.inputs.json"],"engineOptions": "--no-cache"}`
	testFullWorkflowLocalUrl = testProjectFileDir + "/" + testWorkflowLocalUrl
	testTempDir              = "/directory/workflow"
	testWorkflowS3Url        = "s3://workflow/path/file.wdl"
	testWorkflowInvalidUrl   = ":NotURL:"
	testCompressedTmpPath    = "/tmp/123/workflow_1343535"
	testArgsFileName         = "args.txt"
	testArgumentsDir         = "workflow/path/"
	testArgumentsPath        = testArgumentsDir + testArgsFileName
	testOptionFileName       = "test.json"
	testOptionFilePath       = "file://path/to/" + testOptionFileName
	testWesUrl               = "https://TestWesUrl.com/prod"
	testContext1Stack        = "Agc-Context-TestProject1-" + testUserId + "-TestContext1"
	testContext2Stack        = "Agc-Context-TestProject1-" + testUserId + "-TestContext2"
	testRun1Id               = "TestRun1Id"
	testRun2Id               = "TestRun2Id"
	testStackId              = "TestStackId"
	testUserId               = "bender123"
	testTmpAttachmentPath    = "/tmp/attachment.file"
	testProjectFileDir       = "/tmp/path/to/project"
	testErrorPrefix          = "unable to run workflow: "
	testFilePathKey          = "project/" + testProjectName + "/userid/" + testUserId + "/data"
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
	mockInputClient   *storagemocks.MockInputClient
	mockOs            *iomocks.MockOS
	mockZip           *iomocks.MockZip
	mockTmp           *iomocks.MockTmp
	mockFileInfo      *iomocks.MockFileInfo
	mockWes           *wesmocks.MockWesClient
	origRemoveFile    func(name string) error
	origCompressToTmp func(srcPath string) (string, error)
	origWriteToTmp    func(namePattern, content string) (string, error)

	testProjSpec         spec.Project
	wfInstance           ddb.WorkflowInstance
	testStackInfo        cfn.StackInfo
	workAbsDir           string
	inputsAbsDir         string
	testAppendedMANIFEST string

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
	s.mockInputClient = storagemocks.NewMockInputClient(s.ctrl)
	s.mockOs = iomocks.NewMockOS(s.ctrl)
	s.mockZip = iomocks.NewMockZip(s.ctrl)
	s.mockTmp = iomocks.NewMockTmp(s.ctrl)
	s.mockFileInfo = iomocks.NewMockFileInfo(s.ctrl)
	s.mockWes = wesmocks.NewMockWesClient(s.ctrl)

	s.origRemoveFile, removeFile, removeAll = removeFile, s.mockOs.Remove, s.mockOs.RemoveAll
	s.origCompressToTmp, compressToTmp = compressToTmp, s.mockZip.CompressToTmp
	s.origWriteToTmp, writeToTmp, createTempDir = writeToTmp, s.mockTmp.Write, s.mockTmp.TempDir
	osStat = s.mockOs.Stat
	copyFileRecursivelyToLocation = func(destinationDir string, sourceDir string) error {
		return nil
	}

	var err error
	s.workAbsDir, err = os.Getwd()
	require.NoError(s.T(), err)
	s.inputsAbsDir = filepath.Join(s.workAbsDir, testArgumentsDir)
	s.testAppendedMANIFEST = "{\"mainWorkFlowURL\":\"haplotypecaller-gvcf-gatk4.wdl\",\"inputFileURLs\":[\"haplotypecaller-gvcf-gatk4.hg38.wgs.inputs.json\",\"" + filepath.Join(testArgumentsDir, testArgsFileName) + "\"],\"engineOptions\":\"--no-cache\"}"

	s.manager = &Manager{
		Project:     s.mockProjectClient,
		Ssm:         s.mockSsmClient,
		Cfn:         s.mockCfn,
		S3:          s.mockS3Client,
		Ddb:         s.mockDdb,
		Storage:     s.mockStorageClient,
		Config:      s.mockConfigClient,
		InputClient: s.mockInputClient,
		WesFactory:  func(_ string) (wes.Interface, error) { return s.mockWes, nil },
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
	s.mockProjectClient.EXPECT().GetLocation().AnyTimes().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").AnyTimes().Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	uploadCall := s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputS3), nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testMANIFESTPath).Return([]byte(testMANIFEST), nil)
	eq := gomock.GotFormatterAdapter(
		gomock.GotFormatterFunc(
			func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			}),
		gomock.WantFormatter(
			gomock.StringerFunc(func() string {
				return fmt.Sprintf("%s", s.testAppendedMANIFEST)
			}),
			gomock.Eq([]byte(s.testAppendedMANIFEST)),
		),
	)
	s.mockStorageClient.EXPECT().WriteFromBytes(testMANIFESTPath, eq).Return(nil)
	s.mockTmp.EXPECT().Write(testArgsFileName+"_*", testInputS3).Return(testTmpAttachmentPath, nil)
	testInputS3Map := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testInputS3), &testInputS3Map)
	s.mockInputClient.EXPECT().UpdateInputs(s.inputsAbsDir, testInputS3Map, testOutputBucket, testFilePathKey).Return(testInputS3Map, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).After(uploadCall).Return(nil)
	s.mockOs.EXPECT().Remove(testTmpAttachmentPath).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath, "")
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_LocalFile_WithLocalArgs() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().AnyTimes().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").AnyTimes().Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	uploadCall := s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputLocal), nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testMANIFESTPath).Return([]byte(testMANIFEST), nil)
	eq := gomock.GotFormatterAdapter(
		gomock.GotFormatterFunc(
			func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			}),
		gomock.WantFormatter(
			gomock.StringerFunc(func() string {
				return fmt.Sprintf("%s", s.testAppendedMANIFEST)
			}),
			gomock.Eq([]byte(s.testAppendedMANIFEST)),
		),
	)
	s.mockStorageClient.EXPECT().WriteFromBytes(testMANIFESTPath, eq).Return(nil)
	s.mockTmp.EXPECT().Write(testArgsFileName+"_*", testInputLocalToS3).Return(testTmpAttachmentPath, nil)
	testInputS3Map := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testInputLocal), &testInputS3Map)
	testOutputS3Map := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testInputLocalToS3), &testOutputS3Map)
	s.mockInputClient.EXPECT().UpdateInputs(s.inputsAbsDir, testInputS3Map, testOutputBucket, testFilePathKey).Return(testOutputS3Map, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).After(uploadCall).Return(nil)
	s.mockOs.EXPECT().Remove(testTmpAttachmentPath).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath, "")
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_LocalFile_NoArgs() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	uploadCall := s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).After(uploadCall).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_LocalFile_OptionsFile() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testOptionFilePath).Return([]byte(testOptionFileLocal), nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.wfInstance.WorkflowName = testS3WorkflowName
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	actualId, err := s.manager.RunWorkflow(testContext1Name, testS3WorkflowName, "", testOptionFilePath)
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_S3Object_WithLocalArgs() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	fmt.Printf("%s, Jonathan", testArgumentsPath)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputLocal), nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testMANIFESTPath).Return([]byte(testMANIFEST), nil)
	eq := gomock.GotFormatterAdapter(
		gomock.GotFormatterFunc(
			func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			}),
		gomock.WantFormatter(
			gomock.StringerFunc(func() string {
				return fmt.Sprintf("%s", s.testAppendedMANIFEST)
			}),
			gomock.Eq([]byte(s.testAppendedMANIFEST)),
		),
	)
	s.mockStorageClient.EXPECT().WriteFromBytes(testMANIFESTPath, eq).Return(nil)
	s.mockTmp.EXPECT().Write(testArgsFileName+"_*", testInputLocalToS3).Return(testTmpAttachmentPath, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockProjectClient.EXPECT().GetLocation().AnyTimes().Return(testProjectFileDir)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").AnyTimes().Return(testTempDir, nil)
	testInputS3Map := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testInputLocal), &testInputS3Map)
	testOutputS3Map := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testInputLocalToS3), &testOutputS3Map)
	s.mockInputClient.EXPECT().UpdateInputs(s.inputsAbsDir, testInputS3Map, testOutputBucket, testFilePathKey).Return(testOutputS3Map, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.wfInstance.WorkflowName = testS3WorkflowName
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().Remove(testTmpAttachmentPath).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testS3WorkflowName, testArgumentsPath, "")
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
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testS3WorkflowName, "", "")
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_ReadProjectSpecFailure() {
	errorMessage := "failed to read project specification"
	s.mockProjectClient.EXPECT().Read().Return(spec.Project{}, errors.New(errorMessage))
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, "", "", "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_MissingWorkflowSpec() {
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, "dummy", "", "")
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
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testInvalidWorkflowName, "", "")
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
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(false)
	s.mockZip.EXPECT().CompressToTmp(testFullWorkflowLocalUrl).Return("", errors.New(errorMessage))
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
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
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
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
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
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
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte{}, errors.New(errorMessage))
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath, "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_UploadInputFailed() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot upload input"
	expectedInfix := "unable to sync s3://TestOutputBucket/project/TestProject1/userid/bender123/data: "
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().AnyTimes().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").AnyTimes().Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputLocal), nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testMANIFESTPath).Return([]byte(testMANIFEST), nil)
	eq := gomock.GotFormatterAdapter(
		gomock.GotFormatterFunc(
			func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			}),
		gomock.WantFormatter(
			gomock.StringerFunc(func() string {
				return fmt.Sprintf("%s", s.testAppendedMANIFEST)
			}),
			gomock.Eq([]byte(s.testAppendedMANIFEST)),
		),
	)
	s.mockStorageClient.EXPECT().WriteFromBytes(testMANIFESTPath, eq).Return(nil)
	testInputS3Map := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testInputLocal), &testInputS3Map)
	s.mockInputClient.EXPECT().UpdateInputs(s.inputsAbsDir, testInputS3Map, testOutputBucket, testFilePathKey).Return(nil, errors.New(errorMessage))

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath, "")
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
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(false)
	s.mockZip.EXPECT().CompressToTmp(testFullWorkflowLocalUrl).Return(testCompressedTmpPath, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{}, errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
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
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(false)
	s.mockZip.EXPECT().CompressToTmp(testFullWorkflowLocalUrl).Return(testCompressedTmpPath, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	stackInfo := cfn.StackInfo{
		Id:      testStackId,
		Outputs: map[string]string{},
	}
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(stackInfo, nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
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
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return("", errors.New(errorMessage))
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
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
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
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
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_CreateTempDir() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot dir error"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").Return("", errors.New(errorMessage))
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_CopyError() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	errorMessage := "cannot dir error"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").Return(testTempDir, nil)
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(errors.New(errorMessage))
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, "", "")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, testErrorPrefix+errorMessage)
		s.Assert().Empty(actualId)
	}
}

func (s *WorkflowRunTestSuite) TestRunWorkflow_RemoveErrorStillWorks() {
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockProjectClient.EXPECT().GetLocation().AnyTimes().Return(testProjectFileDir)
	s.mockSsmClient.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
	s.mockInputClient.EXPECT().UpdateInputReferencesAndUploadToS3(testFullWorkflowLocalUrl, testTempDir, testOutputBucket, testWorkflowKey).Return(nil)
	s.mockTmp.EXPECT().TempDir("", "workflow_*").AnyTimes().Return(testTempDir, nil)
	s.mockOs.EXPECT().RemoveAll(testTempDir).Return(errors.New("some error"))
	s.mockOs.EXPECT().Stat(testFullWorkflowLocalUrl).Return(s.mockFileInfo, nil)
	s.mockFileInfo.EXPECT().IsDir().Return(true)
	s.mockZip.EXPECT().CompressToTmp(testTempDir).Return(testCompressedTmpPath, nil)
	uploadCall := s.mockS3Client.EXPECT().UploadFile(testOutputBucket, testWorkflowZipKey, testCompressedTmpPath).Return(nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testArgumentsPath).Return([]byte(testInputS3), nil)
	s.mockStorageClient.EXPECT().ReadAsBytes(testMANIFESTPath).Return([]byte(testMANIFEST), nil)
	eq := gomock.GotFormatterAdapter(
		gomock.GotFormatterFunc(
			func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			}),
		gomock.WantFormatter(
			gomock.StringerFunc(func() string {
				return fmt.Sprintf("%s", s.testAppendedMANIFEST)
			}),
			gomock.Eq([]byte(s.testAppendedMANIFEST)),
		),
	)
	s.mockStorageClient.EXPECT().WriteFromBytes(testMANIFESTPath, eq).Return(nil)
	s.mockTmp.EXPECT().Write(testArgsFileName+"_*", testInputS3).Return(testTmpAttachmentPath, nil)
	testInputS3Map := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testInputS3), &testInputS3Map)
	s.mockInputClient.EXPECT().UpdateInputs(s.inputsAbsDir, testInputS3Map, testOutputBucket, testFilePathKey).Return(testInputS3Map, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().RunWorkflow(context.Background(), gomock.Any()).Return(testRun1Id, nil)
	s.mockDdb.EXPECT().WriteWorkflowInstance(context.Background(), s.wfInstance).Return(nil)
	s.mockOs.EXPECT().RemoveAll("extra").Return(nil)
	s.mockOs.EXPECT().Remove(testCompressedTmpPath).After(uploadCall).Return(nil)
	s.mockOs.EXPECT().Remove(testTmpAttachmentPath).Return(nil)

	actualId, err := s.manager.RunWorkflow(testContext1Name, testLocalWorkflowName, testArgumentsPath, "")
	if s.Assert().NoError(err) {
		s.Assert().Equal(testRun1Id, actualId)
	}
}

func TestWorkflowRunTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowRunTestSuite))
}
