package wes

import (
	"context"
	"io"

	"github.com/aws/amazon-genomics-cli/internal/pkg/wes/option"
	"github.com/rs/zerolog/log"
	wes "github.com/rsc/wes_client"
)

type Client struct {
	wes apiClient
}

func New(wesBaseUrl string, profile string) (*Client, error) {
	wesApiClient, err := establishWesConnection(wesBaseUrl+"ga4gh/wes/v1", profile)
	if err != nil {
		return nil, err
	}
	return &Client{wes: wesApiClient}, nil
}

func (c *Client) RunWorkflow(ctx context.Context, options ...option.Func) (string, error) {
	params := new(wes.RunWorkflowOpts)
	for _, optionFunc := range options {
		err := optionFunc(params)
		if err != nil {
			return "", err
		}
	}
	runId, _, err := c.wes.RunWorkflow(ctx, params)
	return runId.RunId, err
}

func (c *Client) StopWorkflow(ctx context.Context, id string) error {
	runId, response, err := c.wes.CancelRun(ctx, id)
	if err != nil {
		log.Error().Msgf("Error stopping workflow instance '%s', the workflow engine is unable to find and/or stop the specified instance", id)
		return err
	}
	log.Debug().Msgf("Stopped workflow '%s', https response is '%s'", runId.RunId, response.Status)
	return nil
}

func (c *Client) GetRunStatus(ctx context.Context, runId string) (string, error) {
	runStatus, _, err := c.wes.GetRunStatus(ctx, runId)
	return string(runStatus.State), err
}

func (c *Client) GetRunLog(ctx context.Context, runId string) (wes.RunLog, error) {
	runLog, _, err := c.wes.GetRunLog(ctx, runId)
	return runLog, err
}

func (c *Client) GetRunLogData(ctx context.Context, runId string, dataUrl string) (*io.ReadCloser, error) {
	runLogDataStream, _, err := c.wes.GetRunLogData(ctx, runId, dataUrl)
	return runLogDataStream, err
}
