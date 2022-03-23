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

func TestSummary_IsHeadProcessEngine(t *testing.T) {
	type testScenario struct {
		engineName string
		expect     bool
	}

	scenarios := []testScenario{{
		engineName: "other",
		expect:     false,
	}, {
		engineName: constants.CROMWELL,
		expect:     false,
	}, {
		engineName: constants.SNAKEMAKE,
		expect:     true,
	}, {
		engineName: constants.NEXTFLOW,
		expect:     true,
	}, {
		engineName: constants.MINIWDL,
		expect:     true,
	}}

	for _, scenario := range scenarios {
		engine := spec.Engine{Engine: scenario.engineName}
		summary := Summary{Engines: []spec.Engine{engine}}
		assert.Equal(t, scenario.expect, summary.IsHeadProcessEngine())
	}
}
