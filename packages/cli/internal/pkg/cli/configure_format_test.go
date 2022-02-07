package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	textFormat    = "text"
	tableFormat   = "table"
	jsonFormat    = "json"
	invalidFormat = "csv"
)

func TestFormatContextOpts_Execute(t *testing.T) {
	tests := [3]string{textFormat, tableFormat, jsonFormat}
	for _, format := range tests {
		t.Run(format, func(t *testing.T) {
			mocks := createMocks(t)
			defer mocks.ctrl.Finish()

			formatContextOpts, err := newFormatContextOpts(formatContextVars{
				format: format,
			})
			require.NoError(t, err)

			mocks.configMock.EXPECT().SetFormat(format).Return(nil)
			formatContextOpts.configClient = mocks.configMock
			err = formatContextOpts.Execute()
			require.NoError(t, err)
		})
	}
}

func TestFormatContextOpts_Validate_InvalidFormat(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatContextOpts, err := newFormatContextOpts(formatContextVars{
		format: invalidFormat,
	})
	require.NoError(t, err)

	formatContextOpts.configClient = mocks.configMock
	err = formatContextOpts.Validate([]string{invalidFormat})
	require.Error(t, err)
}

func TestFormatContextOpts_Validate_ValidFormat(t *testing.T) {
	tests := [3]string{textFormat, tableFormat, jsonFormat}
	for _, format := range tests {
		t.Run(format, func(t *testing.T) {
			mocks := createMocks(t)
			defer mocks.ctrl.Finish()

			formatContextOpts, err := newFormatContextOpts(formatContextVars{
				format: format,
			})
			require.NoError(t, err)

			formatContextOpts.configClient = mocks.configMock
			err = formatContextOpts.Validate([]string{format})
			require.NoError(t, err)
		})
	}
}
