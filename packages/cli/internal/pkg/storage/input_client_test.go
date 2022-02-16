package storage

import (
	"errors"
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

func generateManifest() spec.Manifest {
	manifest := spec.Manifest{}
	manifest.InputFileUrls = append(manifest.InputFileUrls, testFile1)

	return manifest
}

func TestInputClientTestSuite(t *testing.T) {
	suite.Run(t, new(InputClientTestSuite))
}
