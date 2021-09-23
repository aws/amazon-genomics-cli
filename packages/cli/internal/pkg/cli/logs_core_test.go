package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_logsSharedOpts_setDefaultEndTimeIfEmpty_NoFlags_DefaultsToOneHourBack(t *testing.T) {
	opts := logsSharedOpts{}
	oldNow := now
	defer func() { now = oldNow }()
	now = mockNow
	opts.setDefaultEndTimeIfEmpty()
	expectedTime := testTime.Add(-time.Hour)
	if assert.NotNil(t, opts.startTime) {
		assert.Equal(t, expectedTime, *opts.startTime)
		assert.Nil(t, opts.endTime)
	}
}
