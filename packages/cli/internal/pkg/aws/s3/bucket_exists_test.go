package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketName = "test-bucket-name"
)

var testErrorMessage = "test-error-message"

func (m *S3Mock) HeadBucket(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*s3.HeadBucketOutput), err
	}
	return nil, err
}

func TestClient_AssertBucketExists_WithExists(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadBucket", context.Background(), &s3.HeadBucketInput{Bucket: aws.String(testBucketName)}).
		Return(nil, nil)
	exists, err := client.BucketExists(testBucketName)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestClient_AssertBucketExists_WithNotExists(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadBucket", context.Background(), &s3.HeadBucketInput{Bucket: aws.String(testBucketName)}).
		Return(nil, &types.NotFound{Message: &testErrorMessage})
	exists, err := client.BucketExists(testBucketName)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestClient_AssertBucketExists_WithForbidden(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadBucket", context.Background(), &s3.HeadBucketInput{Bucket: aws.String(testBucketName)}).
		Return(nil, fmt.Errorf(testErrorMessage))
	exists, err := client.BucketExists(testBucketName)
	assert.Error(t, err, testErrorMessage)
	assert.False(t, exists)
}
