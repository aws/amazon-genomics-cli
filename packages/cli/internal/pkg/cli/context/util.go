package context

import (
	"fmt"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
)

func environmentMapToList(environmentMap map[string]string) []string {
	var environmentList []string
	for key, value := range environmentMap {
		environmentList = append(environmentList, fmt.Sprintf("%s=%s", key, value))
	}
	return environmentList
}

func cdkResultToContextResult(cdkResults []cdk.Result) []ProgressResult {
	var results []ProgressResult
	for _, cdkResult := range cdkResults {
		progressResult := ProgressResult{
			cdkResult.UniqueKey,
			cdkResult.Outputs,
			cdkResult.Err,
		}
		results = append(results, progressResult)
	}

	return results
}
