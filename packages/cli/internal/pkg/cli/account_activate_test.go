package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/aws/amazon-genomics-cli/internal/pkg/version"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	testAccountBucketName  = "test-account-bucket"
	testAccountRegion      = "test-account-region"
	testAccountId          = "test-account-id"
	testAccountVpcId       = "test-account-vpc-id"
	testImageTag           = "test-image-tag"
	testWesRepository      = "test-wes-repo"
	testCromwellRepository = "test-cromwell-repo"
	testNextflowRepository = "test-nextflow-repo"
	testMiniwdlRepository  = "test-miniwdl-repo"
)

var (
	testImageRefs = map[string]ecr.ImageReference{
		"WES": {
			RegistryId:     testAccountId,
			Region:         testAccountRegion,
			RepositoryName: testWesRepository,
			ImageTag:       testImageTag,
		},
		"CROMWELL": {
			RegistryId:     testAccountId,
			Region:         testAccountRegion,
			RepositoryName: testCromwellRepository,
			ImageTag:       testImageTag,
		},
		"NEXTFLOW": {
			RegistryId:     testAccountId,
			Region:         testAccountRegion,
			RepositoryName: testNextflowRepository,
			ImageTag:       testImageTag,
		},
		"MINIWDL": {
			RegistryId:     testAccountId,
			Region:         testAccountRegion,
			RepositoryName: testMiniwdlRepository,
			ImageTag:       testImageTag,
		},
	}
)

func TestAccountActivateOpts_Execute(t *testing.T) {
	origVerbose := logging.Verbose
	defer func() { logging.Verbose = origVerbose }()
	logging.Verbose = true

	testCases := map[string]struct {
		vpcId       string
		bucketName  string
		setupMocks  func(*testing.T) mockClients
		expectedErr error
	}{
		"generated bucket with no default VPC": {
			bucketName: "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.stsMock.EXPECT().GetAccount().Return(testAccountId, nil)
				mocks.s3Mock.EXPECT().BucketExists("agc-test-account-id-test-account-region").Return(false, nil)
				vars := []string{
					fmt.Sprintf("AGC_BUCKET_NAME=agc-%s-%s", testAccountId, testAccountRegion),
					fmt.Sprintf("CREATE_AGC_BUCKET=%t", true),
					fmt.Sprintf("AGC_PUBLIC_SUBNETS=%t", false),
					fmt.Sprintf("AGC_VERSION=%s", version.Version),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
			vpcId: "",
		},
		"account error": {
			bucketName: "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.stsMock.EXPECT().GetAccount().Return("", fmt.Errorf("some account error"))
				return mocks
			},
			expectedErr: fmt.Errorf("some account error"),
		},
		"new bucket with no default VPC": {
			bucketName: testAccountBucketName,
			vpcId:      "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(false, nil)
				vars := []string{
					fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
					fmt.Sprintf("CREATE_AGC_BUCKET=%t", true),
					fmt.Sprintf("AGC_PUBLIC_SUBNETS=%t", false),
					fmt.Sprintf("AGC_VERSION=%s", version.Version),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
		},
		"existing bucket with no default VPC": {
			bucketName: testAccountBucketName,
			vpcId:      "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(true, nil)
				vars := []string{
					fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
					fmt.Sprintf("CREATE_AGC_BUCKET=%t", false),
					fmt.Sprintf("AGC_PUBLIC_SUBNETS=%t", false),
					fmt.Sprintf("AGC_VERSION=%s", version.Version),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
		},
		"new bucket with custom VPC": {
			bucketName: testAccountBucketName,
			vpcId:      testAccountVpcId,
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(false, nil)
				vars := []string{
					fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
					fmt.Sprintf("CREATE_AGC_BUCKET=%t", true),
					fmt.Sprintf("AGC_PUBLIC_SUBNETS=%t", false),
					fmt.Sprintf("AGC_VERSION=%s", version.Version),
					fmt.Sprintf("VPC_ID=%s", testAccountVpcId),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
		},
		"bucket exists error": {
			bucketName: testAccountBucketName,
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(false, fmt.Errorf("some bucket exists error"))
				return mocks
			},
			expectedErr: fmt.Errorf("some bucket exists error"),
		},
		"deploy error": {
			bucketName: testAccountBucketName,
			vpcId:      "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(true, nil)
				vars := []string{
					fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
					fmt.Sprintf("CREATE_AGC_BUCKET=%t", false),
					fmt.Sprintf("AGC_PUBLIC_SUBNETS=%t", false),
					fmt.Sprintf("AGC_VERSION=%s", version.Version),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(
					gomock.Any(), vars, "activate").Return(nil, fmt.Errorf("some deploy error"))
				return mocks
			},
			expectedErr: fmt.Errorf("some deploy error"),
		},
		"bootstrap error": {
			bucketName: testAccountBucketName,
			vpcId:      "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				vars := []string{
					fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
					fmt.Sprintf("CREATE_AGC_BUCKET=%t", false),
					fmt.Sprintf("AGC_PUBLIC_SUBNETS=%t", false),
					fmt.Sprintf("AGC_VERSION=%s", version.Version),
				}
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(true, nil)
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(nil, fmt.Errorf("some bootstrap error"))
				return mocks
			},
			expectedErr: fmt.Errorf("some bootstrap error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mocks := tc.setupMocks(t)
			defer mocks.ctrl.Finish()
			opts := &accountActivateOpts{
				accountActivateVars: accountActivateVars{
					bucketName: tc.bucketName,
					vpcId:      tc.vpcId,
				},
				stsClient: mocks.stsMock,
				s3Client:  mocks.s3Mock,
				cdkClient: mocks.cdkMock,
				ecrClient: mocks.ecrMock,
				imageRefs: testImageRefs,
				region:    testAccountRegion,
			}

			err := opts.Execute()
			if tc.expectedErr != nil {
				assert.Equal(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
