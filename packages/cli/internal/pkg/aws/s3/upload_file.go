package s3

import (
	"context"
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) UploadFile(bucketName, key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	uploader := manager.NewUploader(c.s3)
	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	return actionable.FindSuggestionForError(err, actionable.AwsErrorMessageToSuggestedActionMap)
}
