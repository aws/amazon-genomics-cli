package awsresources

import (
	"fmt"
	"path"

	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
)

func RenderContextStackName(projectName, contextName, userId string) string {
	return fmt.Sprintf("%s-Context-%s-%s-%s", constants.ProductName, projectName, userId, contextName)
}

func RenderContextStackNameRegexp(projectName, userId string) string {
	return fmt.Sprintf("^%s-Context-%s-%s-([^\\-]+)$", constants.ProductName, projectName, userId)
}

func RenderBucketContextKey(projectName, userId, contextName string, suffix ...string) string {
	args := append([]string{"project", projectName, "userid", userId, "context", contextName}, suffix...)
	return path.Join(args...)
}

func RenderBucketDataKey(projectName, userId string, suffix ...string) string {
	args := append([]string{"project", projectName, "userid", userId, "data"}, suffix...)
	return path.Join(args...)
}

func RenderBootstrapStackName() string {
	return fmt.Sprintf("%s-%s", constants.ProductName, "CDKToolkit")
}
