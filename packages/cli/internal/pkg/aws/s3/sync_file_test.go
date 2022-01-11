package s3

import (
	"context"
	"fmt"
	"io/ioutil"
	nethttp "net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testBucketKey = "test-bucket-key"
)

func (m *S3Mock) HeadObject(ctx context.Context, input *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*s3.HeadObjectOutput), err
	}
	return nil, err
}

func (m *S3Mock) PutObject(ctx context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*s3.PutObjectOutput), err
	}
	return nil, err
}

func TestClient_SyncFile_NoChange(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	_ = os.WriteFile(inputPath, []byte("inputData"), 0644)
	now := time.Now()
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadObject", context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(&s3.HeadObjectOutput{LastModified: &now, ContentLength: 9}, nil)
	err := client.SyncFile(testBucketName, testBucketKey, inputPath)
	assert.NoError(t, err)
}

func TestClient_SyncFile_TimeChange(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	_ = ioutil.WriteFile(inputPath, []byte("inputData"), 0644)
	then := time.Date(1920, time.July, 25, 0, 0, 0, 0, time.UTC)
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadObject", context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(&s3.HeadObjectOutput{LastModified: &then, ContentLength: 9}, nil)
	client.s3.(*S3Mock).On("PutObject", context.Background(), mock.Anything).Return(&s3.PutObjectOutput{}, nil)
	err := client.SyncFile(testBucketName, testBucketKey, inputPath)
	assert.NoError(t, err)
}

func TestClient_SyncFile_SizeChange(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	_ = os.WriteFile(inputPath, []byte("inputData"), 0644)
	fileStat, _ := os.Stat(inputPath)
	modTime := fileStat.ModTime()
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadObject", context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(&s3.HeadObjectOutput{LastModified: &modTime, ContentLength: 8}, nil)
	client.s3.(*S3Mock).On("PutObject", context.Background(), mock.Anything).Return(&s3.PutObjectOutput{}, nil)
	err := client.SyncFile(testBucketName, testBucketKey, inputPath)
	assert.NoError(t, err)
}

func TestClient_SyncFile_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	_ = os.WriteFile(inputPath, []byte("inputData"), 0644)
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadObject", context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, &http.ResponseError{
		ResponseError: &smithyhttp.ResponseError{Response: &smithyhttp.Response{Response: &nethttp.Response{StatusCode: 404}}},
	})
	client.s3.(*S3Mock).On("PutObject", context.Background(), mock.Anything).Return(&s3.PutObjectOutput{}, nil)
	err := client.SyncFile(testBucketName, testBucketKey, inputPath)
	assert.NoError(t, err)
}

func TestClient_SyncFile_HeadError(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	_ = os.WriteFile(inputPath, []byte("inputData"), 0644)
	client := NewMockClient()
	client.s3.(*S3Mock).On("HeadObject", context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, fmt.Errorf("some head object error"))
	err := client.SyncFile(testBucketName, testBucketKey, inputPath)
	assert.Equal(t, fmt.Errorf("some head object error"), err)
}

func TestClient_SyncFile_LocalError(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	client := NewMockClient()
	err := client.SyncFile(testBucketName, testBucketKey, inputPath)
	assert.Error(t, err)
}
