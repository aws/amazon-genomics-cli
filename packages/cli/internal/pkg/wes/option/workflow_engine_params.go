package option

import (
	"encoding/json"

	"github.com/antihax/optional"
	wes "github.com/rsc/wes_client"
)

func WorkflowEngineParams(params map[string]string) Func {
	return func(opts *wes.RunWorkflowOpts) error {
		workflowEngineParamsJsonBytes, err := json.Marshal(params)
		if err != nil {
			return err
		}
		opts.WorkflowEngineParameters = optional.NewString(string(workflowEngineParamsJsonBytes))
		return nil
	}
}
