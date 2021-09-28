package spec

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
)

const LatestVersion = 1

type Project struct {
	Name          string              `yaml:"name"`
	SchemaVersion int                 `yaml:"schemaVersion"`
	Workflows     map[string]Workflow `yaml:"workflows,omitempty"`
	Data          []Data              `yaml:"data,omitempty"`
	Contexts      map[string]Context  `yaml:"contexts,omitempty"`
}

func (projectSpec *Project) GetContext(contextName string) (Context, error) {
	contextSpec, ok := projectSpec.Contexts[contextName]
	if !ok {
		return Context{}, actionable.NewError(
			fmt.Errorf("context '%s' is not defined in Project '%s' specification", contextName, projectSpec.Name),
			"Please add the context to your project spec and deploy it or specify a different context from the command 'agc context list'",
		)
	}

	return contextSpec, nil
}
