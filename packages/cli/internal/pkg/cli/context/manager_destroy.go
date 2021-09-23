package context

import (
	"fmt"
	"path/filepath"

	"github.com/aws/amazon-genomics-cli/common/aws/cdk"
)

const (
	requiredContextPlaceholder = "placeholder"
)

func (m *Manager) Destroy(contextName string, showProgress bool) error {
	m.readProjectSpec()
	m.readConfig()
	m.setContextEnv(contextName)
	m.setContextPlaceholders()
	m.destroyContext(contextName, showProgress)
	return m.err
}

func (m *Manager) setContextPlaceholders() {
	if m.err != nil {
		return
	}

	m.contextEnv.OutputBucketName = requiredContextPlaceholder
	m.contextEnv.ArtifactBucketName = requiredContextPlaceholder
}

func (m *Manager) destroyContext(contextName string, showProgress bool) {
	contextCmd := func() (cdk.ProgressStream, error) {
		return m.Cdk.DestroyApp(filepath.Join(m.homeDir, cdkAppsDirBase, contextDir), m.contextEnv.ToEnvironmentList())
	}
	description := fmt.Sprintf("Destroying resources for context '%s'...", contextName)
	m.executeCdkHelper(contextCmd, description, showProgress)
}
