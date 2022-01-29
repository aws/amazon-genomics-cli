package cdk

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/awsresources"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
)

func (client Client) Bootstrap(appDir string, context []string, executionName string) (ProgressStream, error) {
	cmdArgs := []string{
		"bootstrap",
		"--toolkit-stack-name", awsresources.RenderBootstrapStackName(),
		"--qualifier", constants.AppTagValue,
		"--tags", fmt.Sprintf("%s=%s", constants.AppTagKey, constants.AppTagValue),
		"--profile", client.profile,
	}

	// Add AGC version and custom tags to the bootstrap stack.
	agcVersionKey := "AGC_VERSION="
	customTagsKey := "CUSTOM_TAGS="
	for _, c := range context {
		if strings.HasPrefix(c, agcVersionKey) {
			version := strings.TrimPrefix(c, agcVersionKey)
			cmdArgs = append(cmdArgs, "--tags", fmt.Sprintf("%s=%s", constants.AgcVersionKey, version))
		}
		if strings.HasPrefix(c, customTagsKey) {
			jsonStr := strings.TrimPrefix(c, customTagsKey)
			tagsMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(jsonStr), &tagsMap)
			if err != nil {
				return nil, err
			}
			for k, v := range tagsMap {
				cmdArgs = append(cmdArgs, "--tags", fmt.Sprintf("%s=%s", k, v))
			}
		}
	}

	cmdArgs = appendContextArguments(cmdArgs, context)
	progressStream, err := executeCdkCommand(appDir, cmdArgs, executionName)
	return progressStream, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
