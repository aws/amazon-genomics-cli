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

func (m *S3Mock) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*s3.ListObjectsV2Output), err
	}
	return nil, err
}

func (m *S3Mock) ListObjectVersions(ctx context.Context, input *s3.ListObjectVersionsInput, opts ...func(*s3.Options)) (*s3.ListObjectVersionsOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*s3.ListObjectVersionsOutput), err
	}
	return nil, err
}

func TestClient_EmptyBucket_Success(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("ListObjectsV2", context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(testBucketName),
	}).Return(&s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key: aws.String(testBucketKey),
			},
		},
	}, nil)
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, nil)

	client.s3.(*S3Mock).On("ListObjectVersions", context.Background(), &s3.ListObjectVersionsInput{
		Bucket: aws.String(testBucketName),
	}).Return(&s3.ListObjectVersionsOutput{
		DeleteMarkers: []types.DeleteMarkerEntry{
			{
				Key:       aws.String(testBucketKey),
				VersionId: aws.String(testBucketVersionId),
			},
		},
		Versions: []types.ObjectVersion{
			{
				Key:       aws.String(testBucketKey),
				VersionId: aws.String(testBucketVersionId),
			},
		},
	}, nil)
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket:    aws.String(testBucketName),
		Key:       aws.String(testBucketKey),
		VersionId: aws.String(testBucketVersionId),
	}).Return(nil, nil)

	err := client.EmptyBucket(testBucketName)
	assert.NoError(t, err)
}

func TestClient_EmptyBucket_Failure_1(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("ListObjectsV2", context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(testBucketName),
	}).Return(nil, fmt.Errorf(testErrorMessage))

	err := client.EmptyBucket(testBucketName)
	assert.Error(t, err)
}

func TestClient_EmptyBucket_Failure_2(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("ListObjectsV2", context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(testBucketName),
	}).Return(&s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key: aws.String(testBucketKey),
			},
		},
	}, nil)
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, fmt.Errorf(testErrorMessage))

	err := client.EmptyBucket(testBucketName)
	assert.Error(t, err)
}

func TestClient_EmptyBucket_Failure_3(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("ListObjectsV2", context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(testBucketName),
	}).Return(&s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key: aws.String(testBucketKey),
			},
		},
	}, nil)
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, nil)

	client.s3.(*S3Mock).On("ListObjectVersions", context.Background(), &s3.ListObjectVersionsInput{
		Bucket: aws.String(testBucketName),
	}).Return(nil, fmt.Errorf(testErrorMessage))

	err := client.EmptyBucket(testBucketName)
	assert.Error(t, err)
}

func TestClient_EmptyBucket_Failure_4(t *testing.T) {
	client := NewMockClient()
	client.s3.(*S3Mock).On("ListObjectsV2", context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(testBucketName),
	}).Return(&s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key: aws.String(testBucketKey),
			},
		},
	}, nil)
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(testBucketKey),
	}).Return(nil, nil)

	client.s3.(*S3Mock).On("ListObjectVersions", context.Background(), &s3.ListObjectVersionsInput{
		Bucket: aws.String(testBucketName),
	}).Return(&s3.ListObjectVersionsOutput{
		DeleteMarkers: []types.DeleteMarkerEntry{
			{
				Key:       aws.String(testBucketKey),
				VersionId: aws.String(testBucketVersionId),
			},
		},
		Versions: []types.ObjectVersion{
			{
				Key:       aws.String(testBucketKey),
				VersionId: aws.String(testBucketVersionId),
			},
		},
	}, nil)
	client.s3.(*S3Mock).On("DeleteObject", context.Background(), &s3.DeleteObjectInput{
		Bucket:    aws.String(testBucketName),
		Key:       aws.String(testBucketKey),
		VersionId: aws.String(testBucketVersionId),
	}).Return(nil, fmt.Errorf(testErrorMessage))

	err := client.EmptyBucket(testBucketName)
	assert.Error(t, err)
}
