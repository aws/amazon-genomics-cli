package workflow

func (m *Manager) ListWorkflows() (map[string]Summary, error) {
	m.readProjectSpec()
	m.readConfig()
	m.initWorkflows()
	m.readWorkflowsFromSpec()
	return m.workflows, m.err
}
