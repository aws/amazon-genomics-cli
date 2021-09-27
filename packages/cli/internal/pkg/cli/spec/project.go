package spec

import "fmt"

const LatestVersion = 1

type Project struct {
	Name          string              `yaml:"name"`
	SchemaVersion int                 `yaml:"schemaVersion"`
	Workflows     map[string]Workflow `yaml:"workflows,omitempty"`
	Data          []Data              `yaml:"data,omitempty"`
	Contexts      map[string]Context  `yaml:"contexts,omitempty"`
}

func GetContext(projectSpec Project, contextName string) (Context, error) {
	contextSpec, ok := projectSpec.Contexts[contextName]
	if !ok {
		return Context{}, fmt.Errorf("context '%s' is not defined in Project '%s' specification", contextName, projectSpec.Name)
	}

	return contextSpec, nil
}
