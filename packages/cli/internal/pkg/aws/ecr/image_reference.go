package ecr

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/rs/zerolog/log"
)

type ImageReference struct {
	RegistryId     string
	Region         string
	RepositoryName string
	ImageTag       string
}

func (c *Client) VerifyImageExists(reference ImageReference) error {
	if !reference.isPopulated() {
		return fmt.Errorf("all fields of an ImageReference must be populated, recieved: '%+v'", reference)
	}

	log.Debug().Msgf("verifying presence of '%s:%s' in region: '%s' of registry (account): '%s'",
		reference.RepositoryName, reference.ImageTag, reference.Region, reference.RegistryId)

	ok, err := c.ImageListable(reference.RegistryId, reference.RepositoryName, reference.ImageTag, reference.Region)
	if err != nil {
		return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	if !ok {
		return actionableerror.New(
			fmt.Errorf("cannot verify the presence of container '%s:%s' in region: '%s' of account: '%s'", reference.RepositoryName, reference.ImageTag, reference.Region, reference.RegistryId),
			"Please check your environment variables and permissions",
		)
	}
	return nil
}

func (r ImageReference) isPopulated() bool {
	return r.Region != "" && r.RegistryId != "" && r.ImageTag != "" && r.RepositoryName != ""
}
