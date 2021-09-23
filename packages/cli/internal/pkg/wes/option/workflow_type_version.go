package option

import (
	"github.com/antihax/optional"
	wes "github.com/rsc/wes_client"
)

func WorkflowTypeVersion(version string) Func {
	return func(opts *wes.RunWorkflowOpts) error {
		opts.WorkflowTypeVersion = optional.NewString(version)
		return nil
	}
}
