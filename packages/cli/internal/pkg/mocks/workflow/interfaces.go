package workflowmocks

import "github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"

type WorkflowManager interface {
	workflow.Interface
}
