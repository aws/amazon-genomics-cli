package context

import (
	"regexp"

	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/awsresources"
	"github.com/aws/amazon-genomics-cli/common/aws/cfn"
)

func (m *Manager) StatusList() ([]Instance, error) {
	m.readProjectSpec()
	m.readConfig()
	m.initContexts()
	m.getLocalContexts()
	return m.getAllContexts()
}

func (m *Manager) getAllContexts() ([]Instance, error) {
	if m.err != nil {
		return nil, m.err
	}
	contextStackNameRegexp := regexp.MustCompile(awsresources.RenderContextStackNameRegexp(m.projectSpec.Name, m.userId))
	stacks, err := m.Cfn.ListStacks(contextStackNameRegexp, cfn.ActiveStacksFilter)
	if err != nil {
		return nil, err
	}

	var allContextStatusList []Instance

	for _, stack := range stacks {
		contextName := contextStackNameRegexp.FindStringSubmatch(stack.Name)[1]
		_, isDefinedInProjectFile := m.contexts[contextName]

		allContextStatusList = append(allContextStatusList, Instance{
			ContextName:            contextName,
			ContextStatus:          mapStackToStatus(stack.Status),
			ContextReason:          stack.StatusReason,
			IsDefinedInProjectFile: isDefinedInProjectFile,
		})
	}

	return allContextStatusList, nil
}
