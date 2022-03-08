package wes

import (
	"context"
	"io"

	"github.com/aws/amazon-genomics-cli/internal/pkg/wes/option"
	wes "github.com/rsc/wes_client"
)

type Interface interface {
	RunWorkflow(ctx context.Context, options ...option.Func) (string, error)
	GetRunStatus(ctx context.Context, runId string) (string, error)
	StopWorkflow(ctx context.Context, runId string) error
	GetRunLog(ctx context.Context, runId string) (wes.RunLog, error)
	GetRunLogData(ctx context.Context, runId string, dataUrl string) (*io.ReadCloser, error)
}
