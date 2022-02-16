package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func (m *S3Mock) DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput, opts ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*s3.DeleteBucketOutput), err
	}
	return nil, err
}

func TestClient_DeleteBucket_Success(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("DeleteBucket", context.Background(), &s3.DeleteBucketInput{
		Bucket: aws.String(testBucketName),
	}).Return(nil, nil)
	err := client.DeleteBucket(testBucketName)
	assert.NoError(t, err)
}

func TestClient_DeleteBucket_Failure(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("DeleteBucket", context.Background(), &s3.DeleteBucketInput{
		Bucket: aws.String(testBucketName),
	}).Return(nil, fmt.Errorf(testErrorMessage))
	err := client.DeleteBucket(testBucketName)
	assert.Error(t, err)
}
