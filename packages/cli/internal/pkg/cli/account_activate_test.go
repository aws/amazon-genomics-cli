package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
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
)

var (
	testAccountBaseEnvVars = []string{
		fmt.Sprintf("ECR_WES_ACCOUNT_ID=%s", testAccountId),
		fmt.Sprintf("ECR_WES_REGION=%s", testAccountRegion),
		fmt.Sprintf("ECR_WES_TAG=%s", testImageTag),
		fmt.Sprintf("ECR_WES_REPOSITORY=%s", testWesRepository),

		fmt.Sprintf("ECR_CROMWELL_ACCOUNT_ID=%s", testAccountId),
		fmt.Sprintf("ECR_CROMWELL_REGION=%s", testAccountRegion),
		fmt.Sprintf("ECR_CROMWELL_TAG=%s", testImageTag),
		fmt.Sprintf("ECR_CROMWELL_REPOSITORY=%s", testCromwellRepository),

		fmt.Sprintf("ECR_NEXTFLOW_ACCOUNT_ID=%s", testAccountId),
		fmt.Sprintf("ECR_NEXTFLOW_REGION=%s", testAccountRegion),
		fmt.Sprintf("ECR_NEXTFLOW_TAG=%s", testImageTag),
		fmt.Sprintf("ECR_NEXTFLOW_REPOSITORY=%s", testNextflowRepository),
	}
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
	}
)

func TestAccountActivateOpts_Execute(t *testing.T) {
	t.Setenv("ECR_WES_REGION", testAccountRegion)
	t.Setenv("ECR_WES_TAG", testImageTag)
	t.Setenv("ECR_WES_REPOSITORY", testWesRepository)
	t.Setenv("ECR_CROMWELL_ACCOUNT_ID", testAccountId)
	t.Setenv("ECR_CROMWELL_REGION", testAccountRegion)
	t.Setenv("ECR_CROMWELL_TAG", testImageTag)
	t.Setenv("ECR_CROMWELL_REPOSITORY", testCromwellRepository)

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
				mocks.cdkMock.EXPECT().DeployApp(
					gomock.Any(),
					append([]string{
						fmt.Sprintf("AGC_BUCKET_NAME=agc-%s-%s", testAccountId, testAccountRegion),
						fmt.Sprintf("CREATE_AGC_BUCKET=%t", true),
					}, testAccountBaseEnvVars...)).Return(mocks.progressStream, nil)
				mocks.ecrMock.EXPECT().VerifyImageExists(gomock.Any()).Return(nil).Times(3)
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
				mocks.ecrMock.EXPECT().VerifyImageExists(gomock.Any()).Return(nil).Times(3)
				mocks.cdkMock.EXPECT().DeployApp(
					gomock.Any(),
					append([]string{
						fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
						fmt.Sprintf("CREATE_AGC_BUCKET=%t", true),
					}, testAccountBaseEnvVars...)).Return(mocks.progressStream, nil)
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
				mocks.ecrMock.EXPECT().VerifyImageExists(gomock.Any()).Return(nil).Times(3)
				mocks.cdkMock.EXPECT().DeployApp(
					gomock.Any(),
					append([]string{
						fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
						fmt.Sprintf("CREATE_AGC_BUCKET=%t", false),
					}, testAccountBaseEnvVars...)).Return(mocks.progressStream, nil)
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
				mocks.ecrMock.EXPECT().VerifyImageExists(gomock.Any()).Return(nil).Times(3)
				baseVars := append([]string{
					fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
					fmt.Sprintf("CREATE_AGC_BUCKET=%t", true),
				}, testAccountBaseEnvVars...)
				mocks.cdkMock.EXPECT().DeployApp(
					gomock.Any(),
					append(baseVars, fmt.Sprintf("VPC_ID=%s", testAccountVpcId))).Return(mocks.progressStream, nil)
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
			expectedErr: fmt.Errorf("An error occurred while activating the account. Error was: 'some account error'"),
		},
		"image does not exist error": {
			bucketName: testAccountBucketName,
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(false, nil)
				mocks.ecrMock.EXPECT().VerifyImageExists(gomock.Any()).Return(fmt.Errorf("some image error"))
				return mocks
			},
			expectedErr: fmt.Errorf("An error occurred while activating the account. Error was: 'some image error'"),
		},
		"deploy error": {
			bucketName: testAccountBucketName,
			vpcId:      "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(true, nil)
				mocks.ecrMock.EXPECT().VerifyImageExists(gomock.Any()).Return(nil).Times(3)
				mocks.cdkMock.EXPECT().DeployApp(
					gomock.Any(),
					append([]string{
						fmt.Sprintf("AGC_BUCKET_NAME=%s", testAccountBucketName),
						fmt.Sprintf("CREATE_AGC_BUCKET=%t", false),
					}, testAccountBaseEnvVars...)).Return(nil, fmt.Errorf("some deploy error"))
				return mocks
			},
			expectedErr: fmt.Errorf("An error occurred while activating the account. Error was: 'some deploy error'"),
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
				assert.Error(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
