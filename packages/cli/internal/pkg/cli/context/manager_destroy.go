package context

import (
	"fmt"
	"path/filepath"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
)

const (
	requiredContextPlaceholder = "placeholder"
)

func (m *Manager) Destroy(contexts []string) []ProgressResult {
	m.readProjectInformation()

	progressStreams, contextsWithStreams := m.getStreamsForCdkDestroys(contexts)

	description := fmt.Sprintf("Destroying resources for context(s) %s", contextsWithStreams)
	m.executeCdkHelper(progressStreams, description)
	return m.progressResults
}

func (m *Manager) destroyAllContexts(contexts []string) {
	if m.err != nil {
		var results []ProgressResult
		for _, context := range contexts {
			results = append(results, ProgressResult{Err: m.err, Context: context})
		}

		m.progressResults = results
		return
	}

	progressStreams, contextsWithStreams := m.getStreamsForCdkDestroys(contexts)

	description := fmt.Sprintf("Deploying resources for context(s) %s", contextsWithStreams)
	m.executeCdkHelper(progressStreams, description)
}

func (m *Manager) getStreamsForCdkDestroys(contexts []string) ([]cdk.ProgressStream, []string) {
	var progressStreams []cdk.ProgressStream
	var contextsWithStreams []string
	for _, contextName := range contexts {
		m.setContextEnv(contextName)
		m.setContextPlaceholders()
		if m.err == nil {
			progressStream := m.destroyContext(contextName)
			if progressStream != nil {
				progressStreams = append(progressStreams, progressStream)
				contextsWithStreams = append(contextsWithStreams, contextName)
			}
		}
		if m.err != nil {
			m.progressResults = append(m.progressResults, ProgressResult{Context: contextName, Err: m.err})
		}
	}

	return progressStreams, contextsWithStreams
}

func (m *Manager) setContextPlaceholders() {
	if m.err != nil {
		return
	}

	m.contextEnv.OutputBucketName = requiredContextPlaceholder
	m.contextEnv.ArtifactBucketName = requiredContextPlaceholder
}

func (m *Manager) destroyContext(contextName string) cdk.ProgressStream {
	contextCmd := func() (cdk.ProgressStream, error) {
		return m.Cdk.DestroyApp(filepath.Join(m.homeDir, cdkAppsDirBase, contextDir), m.contextEnv.ToEnvironmentList(), contextName)
	}

	progressStream, err := contextCmd()
	if err != nil {
		m.progressResults = append(m.progressResults, ProgressResult{Context: contextName, Err: err})
	}
	return progressStream
}
