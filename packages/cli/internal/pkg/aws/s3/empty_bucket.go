package s3

import (
	"context"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (c *Client) EmptyBucket(bucketName string) error {
	if err := c.deleteAllBucketObjects(bucketName); err != nil {
		return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}

	if err := c.deleteAllBucketObjectVersions(bucketName); err != nil {
		return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}

	return nil
}

func (c *Client) deleteAllBucketObjects(bucketName string) error {
	ctx := context.Background()
	input := &s3.ListObjectsV2Input{Bucket: aws.String(bucketName)}
	paginator := s3.NewListObjectsV2Paginator(c.s3, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
		}
		for _, object := range page.Contents {
			key := aws.ToString(object.Key)
			log.Debug().Msgf("Deleting object '%s' in bucket '%s'\n", key, bucketName)
			if err := c.DeleteObject(bucketName, key); err != nil {
				return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
			}
		}
	}
	return nil
}

func (c *Client) deleteAllBucketObjectVersions(bucketName string) error {
	ctx := context.Background()
	input := &s3.ListObjectVersionsInput{Bucket: aws.String(bucketName)}
	for {
		output, err := c.s3.ListObjectVersions(ctx, input)
		if err != nil {
			return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
		}
		for _, marker := range output.DeleteMarkers {
			log.Debug().Msgf("Deleting object delete marker '%s.%s' in bucket '%s'\n", *marker.Key, *marker.VersionId, bucketName)
			if err := c.DeleteObjectVersion(bucketName, aws.ToString(marker.Key), aws.ToString(marker.VersionId)); err != nil {
				return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
			}
		}
		for _, marker := range output.Versions {
			log.Debug().Msgf("Deleting object version marker '%s.%s' in bucket '%s'\n", *marker.Key, *marker.VersionId, bucketName)
			if err := c.DeleteObjectVersion(bucketName, aws.ToString(marker.Key), aws.ToString(marker.VersionId)); err != nil {
				return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
			}
		}

		if output.IsTruncated {
			input.KeyMarker = output.NextKeyMarker
			input.VersionIdMarker = output.NextVersionIdMarker
		} else {
			break
		}
	}
	return nil
}
