package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummary_IsEmpty(t *testing.T) {
	summary := Summary{}
	assert.True(t, summary.IsEmpty())
}

func TestSummary_IsNotEmpty(t *testing.T) {
	summary := Summary{Name: "some-context"}
	assert.False(t, summary.IsEmpty())
}

func TestDetail_IsEmpty(t *testing.T) {
	detail := Detail{}
	assert.True(t, detail.IsEmpty())
}

func TestDetail_IsNotEmpty(t *testing.T) {
	detail := Detail{WesUrl: "amazon.com"}
	assert.False(t, detail.IsEmpty())
}
