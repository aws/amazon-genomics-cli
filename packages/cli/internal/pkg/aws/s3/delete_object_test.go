package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketVersionId = "test-bucket-version-id"
)

func (m *S3Mock) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*s3.DeleteObjectOutput), err
	}
	return nil, err
}

func TestClient_DeleteObject_Success(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, nil)
	err := client.DeleteObject(testBucketName, testBucketKey)
	assert.NoError(t, err)
}

func TestClient_DeleteObject_Failure(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, fmt.Errorf(testErrorMessage))
	err := client.DeleteObject(testBucketName, testBucketKey)
	assert.Error(t, err)
}

func TestClient_DeleteObjectVersion_Success(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket:    aws.String(testBucketName),
		Key:       aws.String(testBucketKey),
		VersionId: aws.String(testBucketVersionId),
	}).Return(nil, nil)
	err := client.DeleteObjectVersion(testBucketName, testBucketKey, testBucketVersionId)
	assert.NoError(t, err)
}

func TestClient_DeleteObjectVersion_Failure(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket:    aws.String(testBucketName),
		Key:       aws.String(testBucketKey),
		VersionId: aws.String(testBucketVersionId),
	}).Return(nil, fmt.Errorf(testErrorMessage))
	err := client.DeleteObjectVersion(testBucketName, testBucketKey, testBucketVersionId)
	assert.Error(t, err)
}
