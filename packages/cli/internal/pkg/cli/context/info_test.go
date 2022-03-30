package context

import (
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
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

func TestSummary_IsServerProcessEngine(t *testing.T) {
	tests := map[string]struct {
		engineName string
		expect     bool
	}{
		"otherIsNotAServer": {
			engineName: "other",
			expect:     false,
		},

		"cromwellIsAServer": {
			engineName: constants.CROMWELL,
			expect:     true,
		},

		"snakeMakeIsNotAServer": {
			engineName: constants.SNAKEMAKE,
			expect:     false,
		},
		"nextFlowIsNotAServer": {
			engineName: constants.NEXTFLOW,
			expect:     false,
		},
		"miniwdlIsNotAServer": {
			engineName: constants.MINIWDL,
			expect:     false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			engine := spec.Engine{Engine: test.engineName}
			summary := Summary{Engines: []spec.Engine{engine}}
			actual := summary.IsServerProcessEngine()
			assert.Equal(t, test.expect, actual)
		})
	}
}
