package cli

import (
	"context"
	"testing"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
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

func Test_splitToBatchesBy(t *testing.T) {
	tests := map[string]struct {
		batchSize int
		strings   []string
		expected  [][]string
	}{
		"nil": {
			batchSize: 123,
			strings:   nil,
			expected:  nil,
		},
		"empty": {
			batchSize: 123,
			strings:   []string{},
			expected:  nil,
		},
		"one not complete batch": {
			batchSize: 123,
			strings:   []string{"foo"},
			expected:  [][]string{{"foo"}},
		},
		"one complete batch": {
			batchSize: 2,
			strings:   []string{"foo", "bar"},
			expected:  [][]string{{"foo", "bar"}},
		},
		"two complete batches": {
			batchSize: 1,
			strings:   []string{"foo", "bar"},
			expected:  [][]string{{"foo"}, {"bar"}},
		},
		"one complete and one not complete": {
			batchSize: 2,
			strings:   []string{"foo1", "bar1", "foo2"},
			expected:  [][]string{{"foo1", "bar1"}, {"foo2"}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := splitToBatchesBy(tt.batchSize, tt.strings)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_fanInChannels_nominal(t *testing.T) {
	sourceFunc := func(ctx context.Context, strings ...string) <-chan cwl.StreamEvent {
		channel := make(chan cwl.StreamEvent)
		go func() {
			defer close(channel)
			for _, s := range strings {
				channel <- cwl.StreamEvent{Logs: []string{s}}
			}
		}()
		return channel
	}

	ctx := context.Background()
	src1 := sourceFunc(ctx, "foo1", "bar1")
	src2 := sourceFunc(ctx, "foo2", "bar2", "something else")
	src3 := sourceFunc(ctx, "singleton3")
	src4 := sourceFunc(ctx)

	combined := fanInChannels(ctx, src1, src2, src3, src4)

	var actual []string
	for event := range combined {
		actual = append(actual, event.Logs[0])
	}
	expected := []string{"foo1", "bar1", "foo2", "bar2", "something else", "singleton3"}
	assert.ElementsMatch(t, actual, expected)
}
