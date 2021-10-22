package cli

import (
	"fmt"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type workflowStatusVars struct {
	MaxInstances int
	InstanceId   string
	WorkflowName string
	ContextName  string
}

type workflowStatusOpts struct {
	workflowStatusVars
	wfManager workflow.StatusManager
}

const workflowMaxInstanceDefault = 20
const workflowMaxAllowedInstance = 1000

func newWorkflowStatusOpts(vars workflowStatusVars) (*workflowStatusOpts, error) {
	return &workflowStatusOpts{
		workflowStatusVars: vars,
		wfManager:          workflow.NewManager(profile),
	}, nil
}

// Validate returns an error if the user's input is invalid.
func (o *workflowStatusOpts) Validate() error {
	if o.MaxInstances <= 0 {
		return fmt.Errorf("max number of workflow instances should be grater than 0, provided value: %d", o.MaxInstances)
	}
	if o.MaxInstances > workflowMaxAllowedInstance {
		return fmt.Errorf("max number of workflow instances should not be greater than 1000, provided value: %d", o.MaxInstances)
	}
	return nil
}

// Execute returns an array of status information records about one or more workflow instances.
func (o *workflowStatusOpts) Execute() ([]types.WorkflowInstance, error) {
	var instanceSummaries []workflow.InstanceSummary
	var err error
	switch {
	case o.InstanceId != "":
		instanceSummaries, err = o.wfManager.StatusWorkflowByInstanceId(o.InstanceId)
	case o.WorkflowName != "":
		instanceSummaries, err = o.wfManager.StatusWorkflowByName(o.WorkflowName, o.MaxInstances)
	case o.ContextName != "":
		instanceSummaries, err = o.wfManager.StatusWorkflowByContext(o.ContextName, o.MaxInstances)
	default:
		instanceSummaries, err = o.wfManager.StatusWorkflowAll(o.MaxInstances)
	}
	if err != nil {
		return nil, err
	}
	workflowInstances := make([]types.WorkflowInstance, len(instanceSummaries))
	for i, instance := range instanceSummaries {
		workflowInstances[i] = types.WorkflowInstance{
			Id:            instance.Id,
			WorkflowName:  instance.WorkflowName,
			ContextName:   instance.ContextName,
			State:         instance.State,
			SubmittedTime: instance.SubmitTime,
			InProject:     instance.InProject,
		}

	}
	return workflowInstances, nil
}

// BuildWorkflowStatusCommand builds the command to show the status information for a specific or for multiple workflow instances in the current project.
func BuildWorkflowStatusCommand() *cobra.Command {
	vars := workflowStatusVars{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status for workflow run(s) in the current project.",
		Args:  cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newWorkflowStatusOpts(vars)
			if err != nil {
				return err
			}
			log.Info().Msgf("Showing workflow run(s). Max Runs: %d", opts.MaxInstances)
			if err := opts.Validate(); err != nil {
				return err
			}
			statuses, err := opts.Execute()
			if err != nil {
				return clierror.New("workflow status", vars, err)
			}
			format.Default.Write(statuses)
			return nil
		}),
	}
	cmd.Flags().IntVar(&vars.MaxInstances, "limit", workflowMaxInstanceDefault, "maximum number of workflow instances to show")
	cmd.Flags().StringVarP(&vars.InstanceId, "run-id", "r", "", "show status of specific workflow run")
	cmd.Flags().StringVarP(&vars.WorkflowName, "workflow-name", "n", "", "show status of workflow runs for a specific workflow name")
	cmd.Flags().StringVarP(&vars.ContextName, "context-name", "c", "", "show status of workflow runs in a specific context")
	return cmd
}
