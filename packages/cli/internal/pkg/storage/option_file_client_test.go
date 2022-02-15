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

type OptionClientTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	mockS3Client   *awsmocks.MockS3Client
	mockFileReader *iomocks.MockFileReader
	mockFileWriter *iomocks.MockFileWriter
	mockJson       *iomocks.MockJson
	mockOs         *iomocks.MockOS
	mockSpec       *iomocks.MockSpec
	optionInstance *OptionInstance
}

func (oc *OptionClientTestSuite) BeforeTest(_, _ string) {
	oc.ctrl = gomock.NewController(oc.T())
	oc.mockS3Client = awsmocks.NewMockS3Client(oc.ctrl)
	oc.mockFileReader = iomocks.NewMockFileReader(oc.ctrl)
	oc.mockFileWriter = iomocks.NewMockFileWriter(oc.ctrl)
	oc.mockJson = iomocks.NewMockJson(oc.ctrl)
	oc.mockOs = iomocks.NewMockOS(oc.ctrl)
	oc.mockSpec = iomocks.NewMockSpec(oc.ctrl)

	ioutilReadFile = oc.mockFileReader.ReadFile
	ioutilWriteFile = oc.mockFileWriter.WriteFile
	jsonUnmarshall = oc.mockJson.Unmarshal
	jsonMarshall = oc.mockJson.Marshal
	specFromJson = oc.mockSpec.FromJson
	osStat = oc.mockOs.Stat

	oc.optionInstance = &OptionInstance{
		S3: oc.mockS3Client,
	}
}

func (oc *OptionClientTestSuite) TestUpdateOptionReferenceAndUploadToS3_ManifestMissing() {
	oc.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, errors.New("some error"))
	err := oc.optionInstance.UpdateOptionReferenceAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	oc.Assert().NoError(err)
}

func (ic *OptionClientTestSuite) TestUpdateOptionReferenceAndUploadToS3_ReadManifestFails() {
	expectedErr := "some error"
	ic.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, nil)
	ic.mockSpec.EXPECT().FromJson(expectedStatDirectory).Return(spec.Manifest{}, errors.New(expectedErr))
	err := ic.optionInstance.UpdateOptionReferenceAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	ic.Assert().Error(err, expectedErr)
}

func (oc *OptionClientTestSuite) TestUpdateOptionReferenceAndUploadToS3_ReadReferenceFails() {
	expectedErr := "some error"
	manifest := generateManifestWithOptionFile()
	oc.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, nil)
	oc.mockSpec.EXPECT().FromJson(expectedStatDirectory).Return(manifest, nil)
	oc.mockFileReader.EXPECT().ReadFile(tempProjectDirectory+"/"+testFile1).Return(nil, errors.New(expectedErr))
	err := oc.optionInstance.UpdateOptionReferenceAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	oc.Assert().Error(err, expectedErr)
}

func (oc *OptionClientTestSuite) TestUpdateOptionReferenceAndUploadToS3_UnmarshallFails() {
	expectedErr := "some error"
	manifest := generateManifestWithOptionFile()
	oc.mockOs.EXPECT().Stat(expectedStatDirectory).Return(nil, nil)
	oc.mockSpec.EXPECT().FromJson(expectedStatDirectory).Return(manifest, nil)
	oc.mockFileReader.EXPECT().ReadFile(tempProjectDirectory+"/"+testFile1).Return(testFile1Bytes, nil)
	oc.mockJson.EXPECT().Unmarshal(testFile1Bytes, gomock.Any()).Return(errors.New(expectedErr))
	err := oc.optionInstance.UpdateOptionReferenceAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)

	oc.Assert().Error(err, expectedErr)
}

func generateManifestWithOptionFile() spec.Manifest {
	manifest := spec.Manifest{}
	manifest.OptionFileUrl = testFile1

	return manifest
}

func TestOptionClientTestSuite(t *testing.T) {
	suite.Run(t, new(OptionClientTestSuite))
}
