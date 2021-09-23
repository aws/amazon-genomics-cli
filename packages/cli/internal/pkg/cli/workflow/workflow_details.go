package workflow

import "github.com/rsc/wes_client"

type Details struct {
	Name         string
	TypeLanguage string
	TypeVersion  string
	Source       string
}
type InstanceSummary struct {
	Id           string
	WorkflowName string
	ContextName  string
	SubmitTime   string
	State        string
	InProject    bool
}

func (i *InstanceSummary) IsInstanceRunning() bool {
	return i.State == string(wes_client.RUNNING) || i.State == string(wes_client.INITIALIZING)
}
