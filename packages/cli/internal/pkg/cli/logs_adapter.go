// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type logsAdapterVars struct {
	logsSharedVars
}

type logsAdapterOpts struct {
	logsAdapterVars
	logsSharedOpts
}

func newLogsAdapterOpts(vars logsAdapterVars) (*logsAdapterOpts, error) {
	return &logsAdapterOpts{
		logsAdapterVars: vars,
		logsSharedOpts: logsSharedOpts{
			ctxManager: context.NewManager(profile),
			cwlClient:  aws.CwlClient(profile),
		},
	}, nil
}

func (o *logsAdapterOpts) Validate() error {
	if err := o.validateFlags(); err != nil {
		return err
	}

	if o.ctxManager != nil {
		ctxMap, err := o.ctxManager.List()
		if err != nil {
			return err
		}

		summary := ctxMap[o.contextName]
		engine := summary.Engines[0].Engine
		if engine == constants.TOIL {
			// The Toil engine doesn't use an adapter, so we shouldn't let the user
			// ask for adapter logs.
			return fmt.Errorf("Context does not use an adapter because it is using the Toil engine")
		}
	}

	return o.parseTime(o.logsSharedVars)
}

func (o *logsAdapterOpts) Execute() error {
	contextInfo, err := o.ctxManager.Info(o.contextName)
	if err != nil {
		return err
	}

	logGroupName := contextInfo.WesLogGroupName
	if o.tail {
		err = o.followLogGroup(logGroupName)
	} else {
		err = o.displayLogGroup(logGroupName, o.startTime, o.endTime, o.filter)
	}

	return err
}

func BuildLogsAdapterCommand() *cobra.Command {
	vars := logsAdapterVars{}
	cmd := &cobra.Command{
		Use:   "adapter -c context_name [-f filter] [-s start_date] [-e end_date] [-l look_back] [-t]",
		Short: "Show workflow adapter logs for a given context.",
		Long: `Show workflow adapter logs for a given context.
If no start, end, or look back periods are set, this command will show logs from the last hour.`,
		Example: `
/code agc logs adapter -c myCtx -s 2021/3/31 -e 2021/4/1 -f ERROR`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newLogsAdapterOpts(vars)
			if err != nil {
				return err
			}
			if err = opts.Validate(); err != nil {
				return err
			}
			opts.setDefaultEndTimeIfEmpty()
			log.Info().Msgf("Showing adapter logs for '%s'", opts.contextName)
			if err = opts.Execute(); err != nil {
				return clierror.New("logs adapter", vars, err)

			}
			return nil
		}),
	}
	vars.setFilterFlags(cmd)
	vars.setContextFlag(cmd)
	return cmd
}
