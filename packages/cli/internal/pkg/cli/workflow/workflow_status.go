package workflow

type StatusManager interface {
	StatusWorkflowAll(numInstances int) ([]InstanceSummary, error)
	StatusWorkflowByInstanceId(instanceId string) ([]InstanceSummary, error)
	StatusWorkflowByName(workflowName string, numInstances int) ([]InstanceSummary, error)
	StatusWorkflowByContext(contextName string, numInstances int) ([]InstanceSummary, error)
}

type TasksManager interface {
	GetRunLog(runId string) (RunLog, error)
	GetWorkflowTasks(runId string) ([]Task, error)
	StatusWorkflowByName(workflowName string, numInstances int) ([]InstanceSummary, error)
}

func (m *Manager) StatusWorkflowAll(numInstances int) ([]InstanceSummary, error) {
	m.readProjectSpec()
	m.readConfig()
	m.loadInstances(numInstances)
	m.populateInstancesState()
	m.setFilteredInstances()
	return m.filteredInstances, m.err
}

func (m *Manager) StatusWorkflowByInstanceId(instanceId string) ([]InstanceSummary, error) {
	m.readProjectSpec()
	m.readConfig()
	m.loadInstance(instanceId)
	m.populateInstancesState()
	m.setFilteredInstances()
	return m.filteredInstances, m.err
}

func (m *Manager) StatusWorkflowByName(workflowName string, numInstances int) ([]InstanceSummary, error) {
	m.readProjectSpec()
	m.readConfig()
	m.loadInstancesByWorkflow(workflowName, numInstances)
	m.populateInstancesState()
	m.setFilteredInstances()
	return m.filteredInstances, m.err
}

func (m *Manager) StatusWorkflowByContext(contextName string, numInstances int) ([]InstanceSummary, error) {
	m.readProjectSpec()
	m.readConfig()
	m.loadInstancesByContext(contextName, numInstances)
	m.populateInstancesState()
	m.setFilteredInstances()
	return m.filteredInstances, m.err
}

func (m *Manager) populateInstancesState() {
	if m.err != nil {
		return
	}
	for contextName, instances := range m.instancesPerContext {
		m.setContext(contextName)
		m.setEngineForWorkflowType(contextName)
		if !m.isContextDeployed(contextName) {
			continue
		}
		m.setContextStackInfo(contextName)
		m.setWesUrl()
		m.setWesClient()
		for _, instance := range instances {
			m.updateInstanceState(instance)
			m.updateInProject(instance)
		}
	}
}
