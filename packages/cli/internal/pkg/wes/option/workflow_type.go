package option

import (
	"strings"

	"github.com/antihax/optional"
	wes "github.com/rsc/wes_client"
)

func WorkflowType(workflowType string) Func {
	return func(opts *wes.RunWorkflowOpts) error {
		opts.WorkflowType = optional.NewString(strings.ToUpper(workflowType))
		return nil
	}
}
