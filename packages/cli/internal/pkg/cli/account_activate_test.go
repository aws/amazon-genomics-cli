package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
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
	testAccountSubnetId1   = "test-account-subnet-id-1"
	testAccountSubnetId2   = "test-account-subnet-id-2"
	testImageTag           = "test-image-tag"
	testWesRepository      = "test-wes-repo"
	testCromwellRepository = "test-cromwell-repo"
	testNextflowRepository = "test-nextflow-repo"
	testMiniwdlRepository  = "test-miniwdl-repo"
	testToilRepository     = "test-toil-repo"
	otherEndpointType      = "OTHER"
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
		"TOIL": {
			RegistryId:     testAccountId,
			Region:         testAccountRegion,
			RepositoryName: testToilRepository,
			ImageTag:       testImageTag,
		},
	}
)

func TestAccountActivateOpts_Execute(t *testing.T) {
	origVerbose := logging.Verbose
	defer func() { logging.Verbose = origVerbose }()
	logging.Verbose = true

	testCases := map[string]struct {
		vpcId        string
		subnets      []string
		bucketName   string
		endpointType string
		setupMocks   func(*testing.T) mockClients
		expectedErr  error
	}{
		"default setup with private endpoint": {
			endpointType: PrivateEndpointType,
			expectedErr:  nil,
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.stsMock.EXPECT().GetAccount().Return(testAccountId, nil)
				mocks.s3Mock.EXPECT().BucketExists("agc-test-account-id-test-account-region").Return(false, nil)
				vars := []string{
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=agc-%s-%s", constants.AgcBucketNameEnvKey, testAccountId, testAccountRegion),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, true),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, PrivateEndpointType),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
		},
		"generated bucket with no default VPC": {
			bucketName: "",
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.stsMock.EXPECT().GetAccount().Return(testAccountId, nil)
				mocks.s3Mock.EXPECT().BucketExists("agc-test-account-id-test-account-region").Return(false, nil)
				vars := []string{
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=agc-%s-%s", constants.AgcBucketNameEnvKey, testAccountId, testAccountRegion),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, true),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, ""),
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
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(false, nil)
				vars := []string{
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcBucketNameEnvKey, testAccountBucketName),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, true),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, ""),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
		},
		"existing bucket with no default VPC": {
			bucketName: testAccountBucketName,
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.s3Mock.EXPECT().BucketExists(testAccountBucketName).Return(true, nil)
				vars := []string{
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcBucketNameEnvKey, testAccountBucketName),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, ""),
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
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcBucketNameEnvKey, testAccountBucketName),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, true),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.VpcIdEnvKey, testAccountVpcId),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
		},
		"new bucket with Custom VPC and specified subnet IDs": {
			vpcId:   testAccountVpcId,
			subnets: []string{testAccountSubnetId1, testAccountSubnetId2},
			setupMocks: func(t *testing.T) mockClients {
				mocks := createMocks(t)
				defer close(mocks.progressStream)
				mocks.stsMock.EXPECT().GetAccount().Return(testAccountId, nil)
				mocks.s3Mock.EXPECT().BucketExists("agc-test-account-id-test-account-region").Return(false, nil)
				vars := []string{
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcBucketNameEnvKey, "agc-test-account-id-test-account-region"),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, true),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.VpcIdEnvKey, testAccountVpcId),
					fmt.Sprintf("%s=%s,%s", constants.AgcVpcSubnetsEnvKey, testAccountSubnetId1, testAccountSubnetId2),
				}
				mocks.cdkMock.EXPECT().Bootstrap(gomock.Any(), vars, "bootstrap").Return(mocks.progressStream, nil)
				mocks.cdkMock.EXPECT().DeployApp(gomock.Any(), vars, "activate").Return(mocks.progressStream, nil)
				return mocks
			},
			expectedErr: nil,
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
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcBucketNameEnvKey, testAccountBucketName),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, ""),
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
					fmt.Sprintf("%s=%t", constants.PublicSubnetsEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcBucketNameEnvKey, testAccountBucketName),
					fmt.Sprintf("%s=%t", constants.CreateBucketEnvKey, false),
					fmt.Sprintf("%s=%s", constants.AgcVersionEnvKey, version.Version),
					fmt.Sprintf("%s=%s", constants.AgcAmiEnvKey, ""),
					fmt.Sprintf("%s=%s", constants.AgcEndpointTypeEnvKey, ""),
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
					bucketName:   tc.bucketName,
					vpcId:        tc.vpcId,
					subnets:      tc.subnets,
					endpointType: tc.endpointType,
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

func Test_accountActivateOpts_validate(t *testing.T) {
	origVerbose := logging.Verbose
	defer func() { logging.Verbose = origVerbose }()
	logging.Verbose = true

	testCases := map[string]struct {
		vpcId         string
		subnets       []string
		publicSubnets bool
		endpointType  string
		expectedErr   error
	}{
		"subnets with VPC ID validates": {
			vpcId:        testAccountVpcId,
			subnets:      []string{testAccountSubnetId1, testAccountSubnetId2},
			endpointType: RegionalEndpointType,
			expectedErr:  nil,
		},
		"subnets without VPC ID is invalid": {
			subnets: []string{testAccountSubnetId1, testAccountSubnetId2},
			expectedErr: &clierror.Error{
				Command:         "account activate",
				CommandVars:     accountActivateVars{subnets: []string{testAccountSubnetId1, testAccountSubnetId2}},
				Cause:           fmt.Errorf("\"subnets\" cannot be supplied without supplying a \"vpc\" ID"),
				SuggestedAction: "use the \"vpc\" flag to supply the identity of the VPC containing the subnets",
			},
		},
		"VPC without subnets is valid": {
			vpcId:        testAccountVpcId,
			endpointType: RegionalEndpointType,
			expectedErr:  nil,
		},
		"Public subnets is valid": {
			publicSubnets: true,
			endpointType:  RegionalEndpointType,
			expectedErr:   nil,
		},
		"Public subnets with specific subnets is invalid": {
			publicSubnets: true,
			endpointType:  RegionalEndpointType,
			subnets:       []string{testAccountSubnetId1},
			expectedErr: &clierror.Error{
				Command:         "account activate",
				CommandVars:     accountActivateVars{publicSubnets: true, subnets: []string{testAccountSubnetId1}, endpointType: RegionalEndpointType},
				Cause:           fmt.Errorf("\"subnets\" cannot be supplied without supplying a \"vpc\" ID"),
				SuggestedAction: "use the \"vpc\" flag to supply the identity of the VPC containing the subnets",
			},
		},
		"Public Subnets with VPC is invalid": {
			publicSubnets: true,
			vpcId:         testAccountVpcId,
			endpointType:  RegionalEndpointType,
			expectedErr: &clierror.Error{
				Command:         "account activate",
				CommandVars:     accountActivateVars{publicSubnets: true, vpcId: testAccountVpcId, endpointType: RegionalEndpointType},
				Cause:           fmt.Errorf("both %[1]q and %[2]q cannot be specified together, as %[2]q involves creating a minimal VPC", accountVpcFlag, publicSubnetsFlag),
				SuggestedAction: "Remove one or both of these flags",
			},
		},
		"Private endpoint type is valid": {
			endpointType: PrivateEndpointType,
			expectedErr:  nil,
		},
		"Other endpoints are not valid": {
			endpointType: otherEndpointType,
			expectedErr: &clierror.Error{
				Command:         "account activate",
				CommandVars:     accountActivateVars{endpointType: otherEndpointType},
				Cause:           fmt.Errorf("invalid endpointType '%s', endpointType must be one of %s or %s", otherEndpointType, RegionalEndpointType, PrivateEndpointType),
				SuggestedAction: "use one of the allowed endpoint types",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			opts := &accountActivateOpts{
				accountActivateVars: accountActivateVars{
					vpcId:         tc.vpcId,
					publicSubnets: tc.publicSubnets,
					subnets:       tc.subnets,
					endpointType:  tc.endpointType,
				},
				imageRefs: testImageRefs,
				region:    testAccountRegion,
			}

			err := opts.validate()
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
