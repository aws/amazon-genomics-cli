package main

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"reflect"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
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
	f := formatVars{format: ""}
	mocks.configMock.EXPECT().GetFormat().Return(defaultFormat, nil)
	origConfigClient := newConfigClient
	var mockConfigClient = func() (*config.Client, error) {
		return &config.Client{
			ConfigInterface: mocks.configMock,
		}, nil
	}
	newConfigClient = mockConfigClient
	defer func() { newConfigClient = origConfigClient }()
	configFormat := setFormatter(f)
	require.True(t, reflect.DeepEqual(configFormat, defaultFormat))
}

func TestSetFormatter_FormatFlagSet(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()
	f := formatVars{format: tableFormat}

	origConfigClient := newConfigClient
	var mockConfigClient = func() (*config.Client, error) {
		return &config.Client{
			ConfigInterface: mocks.configMock,
		}, nil
	}
	newConfigClient = mockConfigClient
	defer func() { newConfigClient = origConfigClient }()
	configFormat := setFormatter(f)
	require.True(t, reflect.DeepEqual(configFormat, tableFormat))
}

func TestValidateFormat_ValidFormat(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()
	f := formatVars{format: tableFormat}

	origConfigClient := newConfigClient
	var mockConfigClient = func() (*config.Client, error) {
		return &config.Client{
			ConfigInterface: mocks.configMock,
		}, nil
	}
	newConfigClient = mockConfigClient
	defer func() { newConfigClient = origConfigClient }()
	testFormat := format.FormatterType(f.format)
	err := ValidateFormat(testFormat)
	require.NoError(t, err)
}

func TestValidateFormat_InvalidFormat(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()
	f := formatVars{format: invalidFormat}

	origConfigClient := newConfigClient
	var mockConfigClient = func() (*config.Client, error) {
		return &config.Client{
			ConfigInterface: mocks.configMock,
		}, nil
	}
	newConfigClient = mockConfigClient
	defer func() { newConfigClient = origConfigClient }()
	testFormat := format.FormatterType(f.format)
	err := ValidateFormat(testFormat)
	require.Error(t, err)
}
