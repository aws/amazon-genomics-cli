package workflow

type OutputManager interface {
	OutputByInstanceId(instanceId string) (map[string]interface{}, error)
}

func (m *Manager) OutputByInstanceId(instanceId string) (map[string]interface{}, error) {
	m.readProjectSpec()
	m.readConfig()
	m.loadInstance(instanceId)
	m.setInstanceSummary()
	m.setContext(m.instanceSummary.ContextName)
	m.setContextStackInfo(m.instanceSummary.ContextName)
	m.setWesUrl()
	m.setWesClient()
	m.setWorkflowRunLogOutputs()

	return m.workflowRunLogOutputs, m.err
}
