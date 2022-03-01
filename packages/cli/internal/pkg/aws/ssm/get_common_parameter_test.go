package ssm

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testParameterSuffix = "test-parameter-suffix"
)

func TestClient_GetCommonParameter_ExpectedValue(t *testing.T) {
	testBucketName := "my-test-bucket"
	mockSsm := new(ssmMockClient)
	client := &Client{mockSsm}
	ctx := context.Background()
	mockSsm.On("GetParameter", ctx, &ssm.GetParameterInput{Name: aws.String("/agc/_common/test-parameter-suffix")}).
		Return(&ssm.GetParameterOutput{Parameter: &types.Parameter{Value: aws.String(testBucketName)}}, nil)

	actual, err := client.GetCommonParameter(testParameterSuffix)
	require.NoError(t, err)
	assert.Equal(t, testBucketName, actual)
}

func TestClient_GetCommonParameter_Error(t *testing.T) {
	mockSsm := new(ssmMockClient)
	client := &Client{mockSsm}
	ctx := context.Background()
	mockSsm.On("GetParameter", ctx, mock.Anything).Return((*ssm.GetParameterOutput)(nil), errors.New(""))

	_, err := client.GetCommonParameter(testParameterSuffix)
	assert.Error(t, err)
}

func TestClient_GetCommonParameter_NoValue(t *testing.T) {
	mockSsm := new(ssmMockClient)
	client := &Client{mockSsm}
	ctx := context.Background()
	mockSsm.On("GetParameter", ctx, mock.Anything).
		Return(&ssm.GetParameterOutput{Parameter: &types.Parameter{Value: nil}}, nil)

	_, err := client.GetCommonParameter(testParameterSuffix)
	assert.Error(t, err)
}

func TestClient_GetOutputBucket_ExpectedValue(t *testing.T) {
	testBucketName := "my-test-bucket"
	mockSsm := new(ssmMockClient)
	client := &Client{mockSsm}
	ctx := context.Background()
	mockSsm.On("GetParameter", ctx, &ssm.GetParameterInput{Name: aws.String("/agc/_common/bucket")}).
		Return(&ssm.GetParameterOutput{Parameter: &types.Parameter{Value: aws.String(testBucketName)}}, nil)

	actual, err := client.GetOutputBucket()
	require.NoError(t, err)
	assert.Equal(t, testBucketName, actual)
}

func TestClient_GetCustomTags_ExpectedValue(t *testing.T) {
	tags := "tags"
	mockSsm := new(ssmMockClient)
	client := &Client{mockSsm}
	ctx := context.Background()
	mockSsm.On("GetParameter", ctx, &ssm.GetParameterInput{Name: aws.String("/agc/_common/customTags")}).
		Return(&ssm.GetParameterOutput{Parameter: &types.Parameter{Value: aws.String(tags)}}, nil)

	actual := client.GetCustomTags()
	assert.Equal(t, tags, actual)
}
