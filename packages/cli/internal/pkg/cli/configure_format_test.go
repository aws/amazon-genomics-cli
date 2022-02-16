package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	textFormat    = "text"
	invalidFormat = "csv"
)

func TestFormatContextOpts_Execute(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatContextOpts, err := newFormatContextOpts(formatContextVars{
		format: textFormat,
	})
	require.NoError(t, err)

	mocks.configMock.EXPECT().SetFormat(textFormat).Return(nil)
	formatContextOpts.configClient = mocks.configMock
	err = formatContextOpts.Execute()
	require.NoError(t, err)
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
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	formatContextOpts, err := newFormatContextOpts(formatContextVars{
		format: textFormat,
	})
	require.NoError(t, err)

	formatContextOpts.configClient = mocks.configMock
	err = formatContextOpts.Validate([]string{textFormat})
	require.NoError(t, err)
}
