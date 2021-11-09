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
	defaultFormat = "text"
	tableFormat   = "table"
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
func TestSetFormatter_FormatFlagNotSet(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatOpts, err := newFormatOpts(formatVars{
		format: "",
	})
	require.NoError(t, err)
	mockedConfig := config.Config{
		Format: config.Format{
			Format: defaultFormat,
		},
	}
	mocks.configMock.EXPECT().GetFormat().Return(mockedConfig.Format.Format, nil)
	formatOpts.configClient = mocks.configMock
	// config format should match the default format
	configFormat := setFormatter(formatOpts)
	require.True(t, reflect.DeepEqual(configFormat, mockedConfig.Format.Format))
}

func TestSetFormatter_FormatFlagSet(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatOpts, err := newFormatOpts(formatVars{
		format: tableFormat,
	})
	require.NoError(t, err)

	formatOpts.configClient = mocks.configMock
	configFormat := setFormatter(formatOpts)
	// config format should match the format flag value
	require.True(t, reflect.DeepEqual(configFormat, formatOpts.formatVars.format))
}
