package workflow

type Interface interface {
	ListWorkflows() (map[string]Summary, error)
	RunWorkflow(contextName, workflowName, argumentsUrl string, optionFileUrl string, engineOptions string) (string, error)
	StatusWorkflowByInstanceId(instanceId string) ([]InstanceSummary, error)
	StatusWorkflowByName(workflowName string, numInstances int) ([]InstanceSummary, error)
	StatusWorkflowByContext(contextName string, numInstances int) ([]InstanceSummary, error)
	StatusWorkflowAll(numInstances int) ([]InstanceSummary, error)
	StopWorkflowInstance(runId string)
	GetWorkflowTasks(runId string) ([]Task, error)
}
