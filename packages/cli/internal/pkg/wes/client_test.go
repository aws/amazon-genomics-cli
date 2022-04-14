package wes

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/antihax/optional"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes/option"
	wes "github.com/rsc/wes_client"
	"github.com/stretchr/testify/assert"
)

type testApi struct {
	t                       *testing.T
	expectedWorkflowRunOpts *wes.RunWorkflowOpts
}

func (api testApi) CancelRun(ctx context.Context, runId string) (wes.RunId, *http.Response, error) {
	return wes.RunId{}, nil, nil
}
func (api testApi) GetRunLog(ctx context.Context, runId string) (wes.RunLog, *http.Response, error) {
	return wes.RunLog{}, nil, nil
}
func (api testApi) GetRunStatus(ctx context.Context, runId string) (wes.RunStatus, *http.Response, error) {
	return wes.RunStatus{}, nil, nil
}
func (api testApi) GetServiceInfo(ctx context.Context) (wes.ServiceInfo, *http.Response, error) {
	return wes.ServiceInfo{}, nil, nil
}
func (api testApi) ListRuns(ctx context.Context, localVarOptionals *wes.ListRunsOpts) (wes.RunListResponse, *http.Response, error) {
	return wes.RunListResponse{}, nil, nil
}
func (api testApi) RunWorkflow(ctx context.Context, localVarOptionals *wes.RunWorkflowOpts) (wes.RunId, *http.Response, error) {
	assert.Equal(api.t, api.expectedWorkflowRunOpts, localVarOptionals)
	return wes.RunId{}, nil, nil
}
func (api testApi) GetRunLogData(ctx context.Context, runId string, dataUrl string) (*io.ReadCloser, *http.Response, error) {
	return nil, nil, nil
}

func TestClient_RunWorkflow(t *testing.T) {
	testCases := map[string]struct {
		setupMocks   func(*testing.T) testApi
		inputOptions []option.Func
		expectedErr  error
	}{
		"sets WorkflowEngineParameters": {
			inputOptions: []option.Func{option.WorkflowEngineParams(map[string]string{"key1": "val1"})},
			setupMocks: func(t *testing.T) testApi {
				return testApi{
					t: t,
					expectedWorkflowRunOpts: &wes.RunWorkflowOpts{
						WorkflowEngineParameters: optional.NewString("{\"key1\":\"val1\"}"),
					},
				}
			},
		},
		"skips WorkflowEngineParameters when nil": {
			inputOptions: []option.Func{option.WorkflowEngineParams(nil)},
			setupMocks: func(t *testing.T) testApi {
				return testApi{
					t: t,
					expectedWorkflowRunOpts: &wes.RunWorkflowOpts{
						WorkflowEngineParameters: optional.String{},
					},
				}
			},
		},
		"skips WorkflowEngineParameters when empty": {
			inputOptions: []option.Func{option.WorkflowEngineParams(map[string]string{})},
			setupMocks: func(t *testing.T) testApi {
				return testApi{
					t: t,
					expectedWorkflowRunOpts: &wes.RunWorkflowOpts{
						WorkflowEngineParameters: optional.String{},
					},
				}
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockApi := tc.setupMocks(t)
			client := &Client{wes: mockApi}

			_, err := client.RunWorkflow(context.Background(), tc.inputOptions...)

			if tc.expectedErr != nil {
				assert.Error(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
