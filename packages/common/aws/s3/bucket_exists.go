package s3

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Client) BucketExists(bucketName string) (bool, error) {
	_, err := c.s3.HeadBucket(context.Background(), &s3.HeadBucketInput{Bucket: aws.String(bucketName)})
	if err != nil {
		var errorType *types.NotFound
		if errors.As(err, &errorType) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
