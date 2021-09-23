package option

import (
	"github.com/antihax/optional"
	wes "github.com/rsc/wes_client"
)

func WorkflowUrl(url string) Func {
	return func(opts *wes.RunWorkflowOpts) error {
		opts.WorkflowUrl = optional.NewString(url)
		return nil
	}
}
