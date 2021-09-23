package cli

import (
	"testing"

	awsmocks "github.com/aws/amazon-genomics-cli/cli/internal/pkg/mocks/aws"
	storagemocks "github.com/aws/amazon-genomics-cli/cli/internal/pkg/mocks/storage"
	"github.com/aws/amazon-genomics-cli/common/aws/cdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockClients struct {
	ctrl           *gomock.Controller
	stsMock        *awsmocks.MockStsClient
	s3Mock         *awsmocks.MockS3Client
	cdkMock        *awsmocks.MockCdkClient
	ecrMock        *awsmocks.MockEcrClient
	cfnMock        *awsmocks.MockCfnClient
	configMock     *storagemocks.MockConfigClient
	progressStream cdk.ProgressStream
}

func createMocks(t *testing.T) mockClients {
	ctrl := gomock.NewController(t)

	return mockClients{
		ctrl:           ctrl,
		cdkMock:        awsmocks.NewMockCdkClient(ctrl),
		s3Mock:         awsmocks.NewMockS3Client(ctrl),
		stsMock:        awsmocks.NewMockStsClient(ctrl),
		ecrMock:        awsmocks.NewMockEcrClient(ctrl),
		cfnMock:        awsmocks.NewMockCfnClient(ctrl),
		configMock:     storagemocks.NewMockConfigClient(ctrl),
		progressStream: make(cdk.ProgressStream),
	}
}

func TestSanitizeProjectName(t *testing.T) {
	projectName := "test-project-name"
	sanitizedName := sanitizeProjectName(projectName)
	assert.Equal(t, "testprojectname", sanitizedName)
}

func TestGenerateBucketName(t *testing.T) {
	accountId := "test-account-id"
	region := "test-region"
	bucketName := generateBucketName(accountId, region)
	assert.Equal(t, "agc-test-account-id-test-region", bucketName)
}
