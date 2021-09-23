package s3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testS3Uri = "s3://bucketName/path/to/resource/data.zip"
	testS3Arn = "arn:aws:s3:::bucketName/path/to/resource/data.zip"
)

func TestIsS3Uri_WithS3Uri(t *testing.T) {
	assert.True(t, IsS3Uri("s3://mybucket/path/data.zip"))
}

func TestIsS3Uri_WithHttpsUri(t *testing.T) {
	assert.False(t, IsS3Uri("https://www.amazon.com/"))
}

func TestIsS3Uri_WithBadUri(t *testing.T) {
	assert.False(t, IsS3Uri(string(rune(0x7f))))
}

func TestUriToArn_Success(t *testing.T) {
	arn, err := UriToArn(testS3Uri)
	assert.NoError(t, err)
	assert.Equal(t, testS3Arn, arn)
}

func TestUriToArn_ParseFailure(t *testing.T) {
	_, err := UriToArn(string(rune(0x7f)))
	assert.Error(t, err)
}

func TestUriToArn_SchemeFailure(t *testing.T) {
	_, err := UriToArn("https://s3.console.aws.amazon.com")
	assert.Error(t, err, fmt.Errorf("expected an S3 URI but got 'https'"))
}

func TestRenderS3Uri_WithKey(t *testing.T) {
	uri := RenderS3Uri(testBucketName, testBucketKey)
	assert.Equal(t, "s3://test-bucket-name/test-bucket-key", uri)
}

func TestRenderS3Uri_WithoutKey(t *testing.T) {
	uri := RenderS3Uri(testBucketName, "")
	assert.Equal(t, "s3://test-bucket-name", uri)
}
