package context

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/util"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/rs/zerolog/log"
)

func (m *Manager) Deploy(contextName string, showProgress bool) error {
	m.readProjectSpec()
	m.readConfig()
	m.readContextSpec(contextName)
	m.setDataBuckets()
	m.setOutputBucket()
	m.setArtifactUrl()
	m.setArtifactBucket()
	m.setTaskContext(contextName)
	m.clearCdkContext(contextDir)
	m.setContextEnv(contextName)
	m.deployContextWithTimeout(contextName, showProgress)
	return m.err
}

func (m *Manager) deployContextWithTimeout(contextName string, showProgress bool) {
	if m.err != nil {
		return
	}
	m.err = util.DeployWithTimeout(func() error {
		m.deployContext(contextName, showProgress)
		return m.err
	}, 30*time.Minute)
}

func (m *Manager) clearCdkContext(appDir string) {
	if m.err != nil {
		return
	}
	m.err = m.Cdk.ClearContext(filepath.Join(m.homeDir, cdkAppsDirBase, appDir))
}

func (m *Manager) deployContext(contextName string, showProgress bool) {
	contextCmd := func() (cdk.ProgressStream, error) {
		return m.Cdk.DeployApp(filepath.Join(m.homeDir, cdkAppsDirBase, contextDir), m.contextEnv.ToEnvironmentList())
	}
	description := fmt.Sprintf("Deploying resources for context '%s'...", contextName)
	m.executeCdkHelper(contextCmd, description, showProgress)
}

func (m *Manager) executeCdkHelper(cmd func() (cdk.ProgressStream, error), description string, showProgress bool) {
	if m.err != nil {
		return
	}
	progressStream, err := cmd()
	if err != nil {
		m.err = err
		return
	}
	if !showProgress || logging.Verbose {
		var lastEvent cdk.ProgressEvent
		for event := range progressStream {
			if event.Err != nil {
				for _, line := range lastEvent.Outputs {
					log.Error().Msg(line)
				}
				m.err = event.Err
				return
			}
			lastEvent = event
		}
	} else {
		m.err = progressStream.DisplayProgress(description)
	}
}
