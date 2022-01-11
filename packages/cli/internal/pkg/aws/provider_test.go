package aws

import (
	"context"
	"reflect"
	"testing"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testProfile1 = "test-profile-1"
	testProfile2 = "test-profile-2"
)

func TestClient_UsesCache(t *testing.T) {
	origLoadConfig := loadConfig
	loadConfig = mockLoadConfig
	defer func() { loadConfig = origLoadConfig }()
	testCases := map[string]struct {
		testFunction func() interface{}
		expectedType string
	}{
		"Cdk": {
			testFunction: func() interface{} { return CdkClient(testProfile1) },
			expectedType: "*cdk.Client",
		},
		"Cfn": {
			testFunction: func() interface{} { return CfnClient(testProfile1) },
			expectedType: "*cfn.Client",
		},
		"Cwl": {
			testFunction: func() interface{} { return CwlClient(testProfile1) },
			expectedType: "*cwl.Client",
		},
		"S3": {
			testFunction: func() interface{} { return S3Client(testProfile1) },
			expectedType: "*s3.Client",
		},
		"Ssm": {
			testFunction: func() interface{} { return SsmClient(testProfile1) },
			expectedType: "*ssm.Client",
		},
		"Sts": {
			testFunction: func() interface{} { return StsClient(testProfile1) },
			expectedType: "*sts.Client",
		},
		"Ddb": {
			testFunction: func() interface{} { return DdbClient(testProfile1) },
			expectedType: "*ddb.Client",
		},
		"Batch": {
			testFunction: func() interface{} { return BatchClient(testProfile1) },
			expectedType: "*batch.Client",
		},
		"Ecr": {
			testFunction: func() interface{} { return EcrClient(testProfile1) },
			expectedType: "*ecr.Client",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client1 := tc.testFunction()
			client2 := tc.testFunction()
			require.NotNil(t, client1)
			require.NotNil(t, client2)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(client1).String())
			assert.Equal(t, tc.expectedType, reflect.TypeOf(client2).String())
			assert.Equal(t, &client1, &client2)
		})
	}
}

func TestClient_DoesNotUseCache(t *testing.T) {
	origLoadConfig := loadConfig
	loadConfig = mockLoadConfig
	defer func() { loadConfig = origLoadConfig }()
	testCases := map[string]struct {
		testFunction func(string) interface{}
		expectedType string
	}{
		"Cdk": {
			testFunction: func(profile string) interface{} { return CdkClient(profile) },
			expectedType: "*cdk.Client",
		},
		"Cfn": {
			testFunction: func(profile string) interface{} { return CfnClient(profile) },
			expectedType: "*cfn.Client",
		},
		"Cwl": {
			testFunction: func(profile string) interface{} { return CwlClient(profile) },
			expectedType: "*cwl.Client",
		},
		"S3": {
			testFunction: func(profile string) interface{} { return S3Client(profile) },
			expectedType: "*s3.Client",
		},
		"Ssm": {
			testFunction: func(profile string) interface{} { return SsmClient(profile) },
			expectedType: "*ssm.Client",
		},
		"Sts": {
			testFunction: func(profile string) interface{} { return StsClient(profile) },
			expectedType: "*sts.Client",
		},
		"Ddb": {
			testFunction: func(profile string) interface{} { return DdbClient(profile) },
			expectedType: "*ddb.Client",
		},
		"Batch": {
			testFunction: func(profile string) interface{} { return BatchClient(profile) },
			expectedType: "*batch.Client",
		},
		"Ecr": {
			testFunction: func(profile string) interface{} { return EcrClient(profile) },
			expectedType: "*ecr.Client",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client1 := tc.testFunction(testProfile1)
			client2 := tc.testFunction(testProfile2)
			require.NotNil(t, client1)
			require.NotNil(t, client2)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(client1).String())
			assert.Equal(t, tc.expectedType, reflect.TypeOf(client2).String())
			assert.NotEqual(t, &client1, &client2)
		})
	}
}

func mockLoadConfig(_ context.Context, _ ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
	return aws.Config{}, nil
}
