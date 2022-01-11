package util

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestTimeToAws_WithTime(t *testing.T) {
	timestampSeconds := int64(766134000)
	someTime := time.Unix(timestampSeconds, 0)
	awsTime := TimeToAws(&someTime)
	assert.Equal(t, timestampSeconds*1000, *awsTime)
}

func TestTimeToAws_WithNil(t *testing.T) {
	awsTime := TimeToAws(nil)
	assert.Nil(t, awsTime)
}

func TestTimeFromAws_WithTime(t *testing.T) {
	timestampMillis := int64(766134000000)
	someTime := time.Unix(timestampMillis/1000, 0)
	awsTime := TimeFromAws(&timestampMillis)
	assert.True(t, someTime.Equal(awsTime))
}

func TestTimeFromAws_WithNil(t *testing.T) {
	awsTime := TimeFromAws(nil)
	assert.Equal(t, time.Unix(0, 0), awsTime)
}
