package version

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Api interface {
	s3.ListObjectsV2APIClient
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Store interface {
	ReadVersions(version string, currentTime time.Time) ([]Info, error)
}
