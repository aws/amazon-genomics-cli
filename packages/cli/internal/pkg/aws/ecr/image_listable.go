package ecr

import (
	"context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func (c *Client) ImageListable(registry string, repositoryName string, imageTag string, repositoryRegion string) (bool, error) {
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repositoryName),
		RegistryId:     aws.String(registry),
	}
	paginator := ecr.NewListImagesPaginator(c.ecr, input)
	for paginator.HasMorePages() {
		listImagesOutput, err := paginator.NextPage(context.Background(), func(options *ecr.Options) {
			options.Region = repositoryRegion
		})
		if err != nil {
			return false, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
		}

		for _, image := range listImagesOutput.ImageIds {
			if image.ImageTag != nil && *image.ImageTag == imageTag {
				return true, nil
			}
		}
	}

	return false, nil
}
