package wes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog/log"
	wes "github.com/rsc/wes_client"
)

const (
	numTimesToRetryStartUpPing    = 3
	numSecondsBetweenStartUpPings = 3
)

type apiClient interface {
	CancelRun(ctx context.Context, runId string) (wes.RunId, *http.Response, error)
	GetRunLog(ctx context.Context, runId string) (wes.RunLog, *http.Response, error)
	GetRunStatus(ctx context.Context, runId string) (wes.RunStatus, *http.Response, error)
	GetServiceInfo(ctx context.Context) (wes.ServiceInfo, *http.Response, error)
	ListRuns(ctx context.Context, localVarOptionals *wes.ListRunsOpts) (wes.RunListResponse, *http.Response, error)
	RunWorkflow(ctx context.Context, localVarOptionals *wes.RunWorkflowOpts) (wes.RunId, *http.Response, error)
	GetRunLogData(ctx context.Context, runId string, dataUrl string) (*io.ReadCloser, *http.Response, error)
}

func establishWesConnection(wesUrl string, profile string) (apiClient, error) {
	log.Debug().Msgf("EstablishWesConnection(%s)", wesUrl)
	apiClient, err := getApiClient(wesUrl, profile)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Ping WES endpoint until we get an answer:")
	success := false
	for i := 0; i < numTimesToRetryStartUpPing; i++ {
		log.Debug().Msgf("Attempt %d of %d", i+1, numTimesToRetryStartUpPing)
		isHealthy, err := workflowEngineInstanceIsHealthy(apiClient)
		if (err != nil) || !isHealthy {
			log.Warn().Msgf("Call to WES endpoint '%s' failed: %v", wesUrl, err)
			time.Sleep(time.Second * time.Duration(numSecondsBetweenStartUpPings))
		} else {
			success = true
			break
		}
	}
	if success {
		log.Debug().Msg("Connected to WES endpoint")
	} else {
		apiClient = nil
		err = fmt.Errorf("attempts to establish a connection to WES service endpoint %s timed out", wesUrl)
	}

	return apiClient, err
}

func getApiClient(url string, profile string) (apiClient, error) {
	configuration := wes.NewConfiguration()
	configuration.BasePath = url
	config, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}
	apiClient, err := wes.NewAPISignedClient(configuration, config)
	if err != nil {
		return nil, err
	}

	return apiClient.WorkflowExecutionServiceApi, nil
}

func workflowEngineInstanceIsHealthy(apiClient apiClient) (bool, error) {
	_, _, err := apiClient.GetServiceInfo(context.Background())
	if err != nil {
		return false, err
	}
	return true, nil
}
