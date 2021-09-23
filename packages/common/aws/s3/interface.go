package s3

import (
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Interface interface {
	BucketExists(string) (bool, error)
	SyncFile(bucketName, key, filePath string) error
	UploadFile(bucketName, key, filePath string) error
}

type s3Interface interface {
	s3.HeadBucketAPIClient
	s3.HeadObjectAPIClient
	manager.UploadAPIClient
}
