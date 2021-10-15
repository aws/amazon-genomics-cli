package managermocks

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
)

type WorkflowManager interface {
	workflow.TasksManager
	workflow.StatusManager
	workflow.OutputManager
}
