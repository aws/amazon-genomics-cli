package option

import (
	"encoding/json"

	"github.com/antihax/optional"
	wes "github.com/rsc/wes_client"
)

func WorkflowParams(params map[string]string) Func {
	return func(opts *wes.RunWorkflowOpts) error {
		workflowParamsJsonBytes, err := json.Marshal(params)
		if err != nil {
			return err
		}
		opts.WorkflowParams = optional.NewString(string(workflowParamsJsonBytes))
		return nil
	}
}
