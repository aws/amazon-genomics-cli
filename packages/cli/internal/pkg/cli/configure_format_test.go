package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	textFormat = "text"
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
