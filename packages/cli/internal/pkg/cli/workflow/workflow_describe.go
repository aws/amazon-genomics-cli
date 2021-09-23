package workflow

func (m *Manager) DescribeWorkflow(workflowName string) (Details, error) {
	m.readProjectSpec()
	m.setWorkflowSpec(workflowName)
	m.readConfig()

	return m.renderWorkflowDetails(workflowName)
}
