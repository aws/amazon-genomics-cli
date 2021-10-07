package workflow

import (
	"context"
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
)

type OutputManager interface {
	OutputByInstanceId(instanceId string) (map[string]interface{}, error)
}

func (m *Manager) OutputByInstanceId(instanceId string) (map[string]interface{}, error) {
	m.readProjectSpec()
	m.readConfig()
	m.loadInstance(instanceId)
	if len(m.instances) < 1 {
		return nil, actionableerror.New(fmt.Errorf("no workflow run found for workflow run '%s'", instanceId), "check the workflow run id and check the workflow was run from the current project")
	}
	instanceSummary := m.instances[0]
	m.setContext(instanceSummary.ContextName)
	m.setContextStackInfo(instanceSummary.ContextName)
	m.setWesUrl()
	m.setWesClient()
	runLog, err := m.wes.GetRunLog(context.Background(), instanceSummary.Id)
	if err != nil {
		return nil, err
	}
	return runLog.Outputs, nil
}
