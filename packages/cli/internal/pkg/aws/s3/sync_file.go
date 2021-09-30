package s3

import (
	"context"
	"errors"
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (c *Client) SyncFile(bucketName, key, filePath string) error {
	shouldSync, err := c.shouldSync(bucketName, key, filePath)
	if err != nil {
		return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	if shouldSync {
		log.Debug().Msgf("Uploading '%s' to '%s'", filePath, RenderS3Uri(bucketName, key))
		return c.UploadFile(bucketName, key, filePath)
	}
	log.Debug().Msgf("Skipping upload for '%s', files already exist at '%s'", filePath, RenderS3Uri(bucketName, key))
	return nil
}

func (c *Client) shouldSync(bucketName, key, filePath string) (bool, error) {
	localFile, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	remoteFile, err := c.getObjectMetadata(bucketName, key)
	if err != nil {
		var responseErr *http.ResponseError
		if errors.As(err, &responseErr) && responseErr.HTTPStatusCode() == 404 {
			return true, nil
		}
		return false, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	differentSize := localFile.Size() != remoteFile.ContentLength
	localNewer := localFile.ModTime().After(*remoteFile.LastModified)
	return differentSize || localNewer, nil
}

func (c *Client) getObjectMetadata(bucketName, key string) (*s3.HeadObjectOutput, error) {
	headObjectOutput, err := c.s3.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	return headObjectOutput, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
