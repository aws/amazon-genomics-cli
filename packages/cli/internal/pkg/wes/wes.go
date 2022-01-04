package wes

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/rs/zerolog/log"
	wes "github.com/rsc/wes_client"
)

const (
	numTimesToRetryStartUpPing    = 3
	numSecondsBetweenStartUpPings = 3
)

func EstablishWesConnection(wesUrl string, profile string) (*wes.APIClient, error) {
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

func getApiClient(url string, profile string) (*wes.APIClient, error) {
	configuration := wes.NewConfiguration()
	configuration.BasePath = url
	config := aws.GetProfileConfig(profile)
	apiClient, err := wes.NewAPISignedClient(configuration, config)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func workflowEngineInstanceIsHealthy(apiClient *wes.APIClient) (bool, error) {
	_, _, err := apiClient.WorkflowExecutionServiceApi.GetServiceInfo(context.Background())
	if err != nil {
		return false, err
	}
	return true, nil
}
