package spec

type Workflow struct {
	Type      WorkflowType `yaml:"type"`
	SourceURL string       `yaml:"sourceURL"`
}

type WorkflowType struct {
	Language string `yaml:"language"`
	Version  string `yaml:"version"`
}
