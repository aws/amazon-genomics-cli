package cli

import (
	"bytes"
	"errors"
	"sort"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/jeremywohl/flatten"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type workflowOutputVars struct {
	runId string
}

type workflowOutputOpts struct {
	vars      workflowOutputVars
	wfManager workflow.OutputManager
}

func newWorkflowOutputOpts(vars workflowOutputVars) (*workflowOutputOpts, error) {
	return &workflowOutputOpts{
		vars:      vars,
		wfManager: workflow.NewManager(profile),
	}, nil
}

func (o *workflowOutputOpts) Validate() error {
	if strings.TrimSpace(o.vars.runId) == "" {
		return actionableerror.New(errors.New("runId contains only white space"), "provide a valid runId")
	}
	return nil
}

func (o *workflowOutputOpts) Execute() ([]types.Output, error) {
	instanceOutput, err := o.wfManager.OutputByInstanceId(o.vars.runId)
	if err != nil {
		return nil, err
	}
	return processOutput(instanceOutput)
}

// BuildWorkflowOutputCommand builds the command to show the output for a workflow instance.
func BuildWorkflowOutputCommand() *cobra.Command {
	vars := workflowOutputVars{}
	cmd := &cobra.Command{
		Use:   "output run_id",
		Short: "Show the output for a workflow run in the current project.",
		Args:  cobra.ExactArgs(1),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			vars.runId = args[0]
			opts, err := newWorkflowOutputOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			log.Info().Msgf("Obtaining final outputs for workflow runId '%s'", vars.runId)
			output, err := opts.Execute()
			if err != nil {
				return clierror.New("workflow output", vars, err)
			}
			format.Default.Write(output)
			return nil
		}),
	}

	return cmd
}

func processOutput(raw map[string]interface{}) ([]types.Output, error) {
	flatOutput, err := flatten.Flatten(raw, "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}

	var output []types.Output
	for key, element := range flatOutput {
		b := bytes.NewBufferString("")
		format.NewStringFormatter(b).Write(element)

		output = append(output, types.Output{
			Path:  key,
			Value: strings.TrimSpace(strings.ReplaceAll(b.String(), "\n", " ")),
		})
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].Path < output[j].Path
	})
	return output, nil
}
