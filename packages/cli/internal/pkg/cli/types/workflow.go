package types

type Workflow struct {
	Name         string
	TypeLanguage string
	TypeVersion  string
	Source       string
}

type WorkflowName struct {
	Name string
}

type WorkflowInstance struct {
	Id            string
	WorkflowName  string
	ContextName   string
	State         string
	SubmittedTime string
	InProject     bool
}

type Output struct {
	Path  string
	Value string
}
