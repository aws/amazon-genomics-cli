package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type InputClientTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	mockS3Client   *awsmocks.MockS3Client
	mockFileReader *iomocks.MockFileReader
	mockFileWriter *iomocks.MockFileWriter
	mockJson       *iomocks.MockJson
	mockOs         *iomocks.MockOS
	mockSpec       *iomocks.MockSpec
	inputInstance  *InputInstance
}

func (ic *InputClientTestSuite) BeforeTest(_, _ string) {
	ic.ctrl = gomock.NewController(ic.T())
	ic.mockS3Client = awsmocks.NewMockS3Client(ic.ctrl)
	ic.mockFileReader = iomocks.NewMockFileReader(ic.ctrl)
	ic.mockFileWriter = iomocks.NewMockFileWriter(ic.ctrl)
	ic.mockJson = iomocks.NewMockJson(ic.ctrl)
	ic.mockOs = iomocks.NewMockOS(ic.ctrl)
	ic.mockSpec = iomocks.NewMockSpec(ic.ctrl)

	ioutilReadFile = ic.mockFileReader.ReadFile
	ioutilWriteFile = ic.mockFileWriter.WriteFile
	jsonUnmarshall = ic.mockJson.Unmarshal
	jsonMarshall = ic.mockJson.Marshal
	specFromJson = ic.mockSpec.FromJson
	osStat = ic.mockOs.Stat
	stat = ic.mockOs.Stat

	ic.inputInstance = &InputInstance{
		S3: ic.mockS3Client,
	}
}

const (
	expectedStatDirectory   = tempProjectDirectory + "/" + ManifestFileName
	initialProjectDirectory = "dir"
	tempProjectDirectory    = "tempDir"
	bucketName              = "s3://bucketName"
	baseS3Key               = "some/key"
	testFile1               = "testFile.json"
	testFile1FullPath       = initialProjectDirectory + "/" + testFile1
)

var (
	testFile1Bytes = []byte("test1")
)

func (ic *InputClientTestSuite) TestUpdateInputReferencesAndUploadToS3_ManifestMissing() {
	ic.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, errors.New("some error"))
	err := ic.inputInstance.UpdateInputReferencesAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	ic.Assert().NoError(err)
}

func (ic *InputClientTestSuite) TestUpdateInputReferencesAndUploadToS3_ReadManifestFails() {
	expectedErr := "some error"
	ic.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, nil)
	ic.mockSpec.EXPECT().FromJson(expectedStatDirectory).Return(spec.Manifest{}, errors.New(expectedErr))
	err := ic.inputInstance.UpdateInputReferencesAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	ic.Assert().Error(err, expectedErr)
}

func (ic *InputClientTestSuite) TestUpdateInputReferencesAndUploadToS3_ReadReferenceFails() {
	expectedErr := "some error"
	manifest := generateManifest()
	ic.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, nil)
	ic.mockSpec.EXPECT().FromJson(expectedStatDirectory).Return(manifest, nil)
	ic.mockFileReader.EXPECT().ReadFile(tempProjectDirectory+"/"+testFile1).Return(nil, errors.New(expectedErr))
	err := ic.inputInstance.UpdateInputReferencesAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	ic.Assert().Error(err, expectedErr)
}

func (ic *InputClientTestSuite) TestUpdateInputReferencesAndUploadToS3_UnmarshallFails() {
	expectedErr := "some error"
	manifest := generateManifest()
	ic.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, nil)
	ic.mockSpec.EXPECT().FromJson(expectedStatDirectory).Return(manifest, nil)
	ic.mockFileReader.EXPECT().ReadFile(tempProjectDirectory+"/"+testFile1).Return(testFile1Bytes, nil)
	ic.mockJson.EXPECT().Unmarshal(testFile1Bytes, gomock.Any()).Return(errors.New(expectedErr))
	err := ic.inputInstance.UpdateInputReferencesAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	ic.Assert().Error(err, expectedErr)
}

func (ic *InputClientTestSuite) TestUpdateInputsInFile_WriteFileFails() {
	inputFile := map[string]interface{}{
		"a": 1,
	}
	inputFileString := []byte("{\"a\":" + testFile1)
	ic.mockJson.EXPECT().Marshal(inputFile).Return(inputFileString, nil)
	expectedErr := errors.New("FileNotFound")
	ic.mockFileWriter.EXPECT().WriteFile(tempProjectDirectory, inputFileString, os.FileMode(0644)).Return(expectedErr)

	err := ic.inputInstance.updateInputsInFile(initialProjectDirectory, inputFile, "bucketName", baseS3Key, tempProjectDirectory)
	ic.Assert().Equal(err, expectedErr)
}

func (ic *InputClientTestSuite) TestUpdateInputsInFile_UploadFileFails() {
	inputFile := map[string]interface{}{
		"a": testFile1,
	}
	mockFileInfo := iomocks.NewMockFileInfo(ic.ctrl)
	mockFileInfo.EXPECT().IsDir().AnyTimes().Return(false)
	ic.mockOs.EXPECT().Stat(testFile1FullPath).AnyTimes().Return(mockFileInfo, nil)
	expectedErr := errors.New("FileNotFound")
	ic.mockS3Client.EXPECT().UploadFile("bucketName", baseS3Key+"/"+testFile1, "dir/"+testFile1).AnyTimes().Return(expectedErr)

	err := ic.inputInstance.updateInputsInFile(initialProjectDirectory, inputFile, "bucketName", baseS3Key, tempProjectDirectory)
	ic.Assert().Equal(err, expectedErr)
}

func (ic *InputClientTestSuite) TestUpdateInputsInFile_MarshallFails() {
	inputFile := map[string]interface{}{
		"a": 1,
	}
	ic.mockOs.EXPECT().Stat(testFile1FullPath).AnyTimes().Return(os.FileInfo(nil), nil)
	expectedErr := errors.New("FileNotFound")
	ic.mockJson.EXPECT().Marshal(inputFile).Return(nil, expectedErr)

	err := ic.inputInstance.updateInputsInFile(initialProjectDirectory, inputFile, "bucketName", baseS3Key, tempProjectDirectory)
	ic.Assert().Equal(err, expectedErr)
}

func (ic *InputClientTestSuite) TestUpdateInputs_HappyCase() {
	inputFile := map[string]interface{}{
		"a": testFile1,
		"b": 1,
		"c": []interface{}{testFile1, 1, "2"},
		"d": []interface{}{[]interface{}{testFile1}},
		"e": "params",
		"f": testFile1 + "," + testFile1,
	}
	expectedUpdatedInputFile := map[string]interface{}{
		"a": "s3://bucketName/some/key/testFile.json",
		"b": 1,
		"c": []interface{}{1, "s3://bucketName/some/key/testFile.json", "2"},
		"d": []interface{}{[]interface{}{testFile1}},
		"e": "params",
		"f": "s3://bucketName/some/key/testFile.json" + "," + "s3://bucketName/some/key/testFile.json",
	}
	mockFileInfo := iomocks.NewMockFileInfo(ic.ctrl)
	mockFileInfo.EXPECT().IsDir().AnyTimes().Return(false)
	ic.mockOs.EXPECT().Stat(testFile1FullPath).AnyTimes().Return(mockFileInfo, nil)
	expectedErr := errors.New("FileNotFound")
	//Using gomock.Any() since there are a bunch of file paths that are being passed around, and this validation is anyway convered in above cases.
	ic.mockOs.EXPECT().Stat(gomock.Any()).AnyTimes().Return(nil, expectedErr)
	ic.mockS3Client.EXPECT().UploadFile("bucketName", gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	actualUpdatedInputFile, err := ic.inputInstance.UpdateInputs(initialProjectDirectory, inputFile, "bucketName", baseS3Key)
	ic.Assert().Equal(err, nil)
	ic.Assert().Equal(actualUpdatedInputFile, expectedUpdatedInputFile)
}

func (ic *InputClientTestSuite) TestUpdateInputs_EmptyString() {
	inputFile := map[string]interface{}{
		"a":"",
		"b":testFile1,
	}
	expectedUpdatedInputFile := map[string]interface{}{
		"a":"",
		"b":"s3://bucketName/some/key/testFile.json",
	}

	mockFileInfo1 := iomocks.NewMockFileInfo(ic.ctrl)
	mockFileInfo1.EXPECT().IsDir().Return(true)
	ic.mockOs.EXPECT().Stat("dir/").Return(mockFileInfo1,nil)

	mockFileInfo2 := iomocks.NewMockFileInfo(ic.ctrl)
	mockFileInfo2.EXPECT().IsDir().Return(false)
	ic.mockOs.EXPECT().Stat("dir/"+testFile1).Return(mockFileInfo2,nil)
	ic.mockS3Client.EXPECT().UploadFile("bucketName", gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	actualUpdatedInputFile, err := ic.inputInstance.UpdateInputs(initialProjectDirectory, inputFile, "bucketName", baseS3Key)
	ic.Assert().NoError(err)
	ic.Assert().Equal(expectedUpdatedInputFile, actualUpdatedInputFile)
}

func generateManifest() spec.Manifest {
	manifest := spec.Manifest{}
	manifest.InputFileUrls = append(manifest.InputFileUrls, testFile1)

	return manifest
}

func TestInputClientTestSuite(t *testing.T) {
	suite.Run(t, new(InputClientTestSuite))
}
