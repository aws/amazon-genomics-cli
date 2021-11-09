package main

import (
	"reflect"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	tableFormat = "table"
)

type mockClients struct {
	ctrl       *gomock.Controller
	configMock *storagemocks.MockConfigClient
}

func createMocks(t *testing.T) mockClients {
	ctrl := gomock.NewController(t)

	return mockClients{
		ctrl:       ctrl,
		configMock: storagemocks.NewMockConfigClient(ctrl),
	}
}
func TestSetFormatter(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatOpts, err := newFormatOpts(formatVars{
		defaultFormat: tableFormat,
	})
	require.NoError(t, err)
	mockedConfig := config.Config{
		Format: config.Format{
			Format: tableFormat,
		},
	}
	formatOpts.configClient = mocks.configMock
	setFormatter(formatOpts)
	// check the actual config to see if it matches the config mock
	configClient, _ := config.NewConfigClient()
	configFormat, _ := configClient.GetFormat()
	require.True(t, reflect.DeepEqual(configFormat, mockedConfig.Format.Format))
}
