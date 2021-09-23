package workflow

func (m *Manager) StopWorkflowInstance(runId string) {
	m.readProjectSpec()
	m.readConfig()
	m.setInstanceToStop(runId)
	m.setContext(m.instanceToStop.ContextName)
	m.setEngineForWorkflowType(m.instanceToStop.ContextName)
	m.setContextStackInfo(m.instanceToStop.ContextName)
	m.setWesUrl()
	m.setWesClient()
	m.stopWorkflowInstance(runId)
}
