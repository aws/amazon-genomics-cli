package context

import (
	"fmt"
	"path/filepath"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
)

func (m *Manager) Deploy(contexts []string) []ProgressResult {
	m.readProjectInformation()
	m.deployAllContexts(contexts)
	return m.progressResults
}

func (m *Manager) deployAllContexts(contexts []string) {
	if m.err != nil {
		var results []ProgressResult
		for _, context := range contexts {
			results = append(results, ProgressResult{Err: m.err, Context: context})
		}

		m.progressResults = results
		return
	}

	progressStreams, contextsWithStreams := m.getStreamsForCdkDeployments(contexts)

	description := fmt.Sprintf("Deploying resources for context(s) %s", contextsWithStreams)
	m.executeCdkHelper(progressStreams, description)
}

func (m *Manager) getStreamsForCdkDeployments(contexts []string) ([]cdk.ProgressStream, []string) {
	var progressStreams []cdk.ProgressStream
	var contextsWithStreams []string
	for _, contextName := range contexts {
		m.readContextSpec(contextName)
		m.setCdkConfigurationForDeployment(contextName)
		m.clearCdkContext(contextDir)
		m.setContextEnv(contextName)

		if m.err == nil {
			progressStream := m.deployContext(contextName)
			if progressStream != nil {
				progressStreams = append(progressStreams, progressStream)
				contextsWithStreams = append(contextsWithStreams, contextName)
			}
		}

		if m.err != nil {
			m.progressResults = append(m.progressResults, ProgressResult{Context: contextName, Err: m.err})
		}
		m.err = nil
	}

	return progressStreams, contextsWithStreams
}

func (m *Manager) setCdkConfigurationForDeployment(contextName string) {
	m.setDataBuckets()
	m.setOutputBucket()
	m.setArtifactUrl()
	m.setArtifactBucket()
	m.setTaskContext(contextName)
}

func (m *Manager) clearCdkContext(appDir string) {
	if m.err != nil {
		return
	}
	m.err = m.Cdk.ClearContext(filepath.Join(m.homeDir, cdkAppsDirBase, appDir))
}

func (m *Manager) deployContext(contextName string) cdk.ProgressStream {
	contextCmd := func() (cdk.ProgressStream, error) {
		return m.Cdk.DeployApp(filepath.Join(m.homeDir, cdkAppsDirBase, contextDir), m.contextEnv.ToEnvironmentList(), contextName)
	}
	progressStream, err := contextCmd()
	if err != nil {
		m.progressResults = append(m.progressResults, ProgressResult{Context: contextName, Err: err})
	}
	return progressStream
}
