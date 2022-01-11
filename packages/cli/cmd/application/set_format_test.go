package main

import (
	"reflect"
	"testing"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
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
	origConfigClient := newConfigClient
	newConfigClient = func() (storage.ConfigClient, error) {
		return mocks.configMock, nil
	}
	mocks.configMock.EXPECT().GetFormat().Return(defaultFormat, nil)
	defer func() { newConfigClient = origConfigClient }()
	configFormat := setFormatter(f)
	require.True(t, reflect.DeepEqual(configFormat, defaultFormat))
}

func TestSetFormatter_FormatFlagSet(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()
	f := formatVars{format: tableFormat}

	origConfigClient := newConfigClient
	newConfigClient = func() (storage.ConfigClient, error) {
		return mocks.configMock, nil
	}
	defer func() { newConfigClient = origConfigClient }()
	configFormat := setFormatter(f)
	require.True(t, reflect.DeepEqual(configFormat, tableFormat))
}

func TestValidateFormat_ValidFormat(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()
	f := formatVars{format: tableFormat}

	origConfigClient := newConfigClient
	newConfigClient = func() (storage.ConfigClient, error) {
		return mocks.configMock, nil
	}
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
	newConfigClient = func() (storage.ConfigClient, error) {
		return mocks.configMock, nil
	}
	defer func() { newConfigClient = origConfigClient }()
	testFormat := format.FormatterType(f.format)
	err := ValidateFormat(testFormat)
	require.Error(t, err)
}
