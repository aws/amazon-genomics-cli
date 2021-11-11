package main

import (
	"reflect"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	tableFormat   = "table"
	invalidFormat = "csv"
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
			Name: defaultFormat,
		},
	}
	mocks.configMock.EXPECT().GetFormat().Return(mockedConfig.Format.Name, nil)
	formatOpts.configClient = mocks.configMock
	configFormat := setFormatter(formatOpts)
	require.True(t, reflect.DeepEqual(configFormat, mockedConfig.Format.Name))
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
	require.True(t, reflect.DeepEqual(configFormat, formatOpts.formatVars.format))
}

func TestValidateFormat_ValidFormat(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatOpts, err := newFormatOpts(formatVars{
		format: tableFormat,
	})
	require.NoError(t, err)

	formatOpts.configClient = mocks.configMock
	format := format.FormatterType(formatOpts.formatVars.format)
	err = ValidateFormat(format)
	require.NoError(t, err)
}

func TestValidateFormat_InvalidFormat(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatOpts, err := newFormatOpts(formatVars{
		format: invalidFormat,
	})
	require.NoError(t, err)

	formatOpts.configClient = mocks.configMock
	format := format.FormatterType(formatOpts.formatVars.format)
	err = ValidateFormat(format)
	require.Error(t, err)
}
