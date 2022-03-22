package context

import (
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/stretchr/testify/assert"
)

var testServerEngineNames = []string{"cromwell"}
var testHeadEngineNames = []string{"nextflow", "miniwdl", "snakemake"}

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

func TestSummary_IsHeadProcessEngine_HeadEnginesShouldReturnTrue(t *testing.T) {
	for _, engineName := range testHeadEngineNames {
		engine := spec.Engine{Engine: engineName}
		summary := Summary{Engines: []spec.Engine{engine}}
		assert.True(t, summary.IsHeadProcessEngine())
	}
}

func TestSummary_IsHeadProcessEngine_ServerEnginesShouldReturnFalse(t *testing.T) {
	for _, engineName := range testServerEngineNames {
		engine := spec.Engine{Engine: engineName}
		summary := Summary{Engines: []spec.Engine{engine}}
		assert.False(t, summary.IsHeadProcessEngine())
	}
}

func TestSummary_IsServerProcessEngine_SeverEnginesShouldReturnTrue(t *testing.T) {
	for _, engineName := range testServerEngineNames {
		engine := spec.Engine{Engine: engineName}
		summary := Summary{Engines: []spec.Engine{engine}}
		assert.True(t, summary.IsServerProcessEngine())
	}
}

func TestSummary_IsServerProcessEngine_HeadProcessEnginesShouldReturnFalse(t *testing.T) {
	for _, engineName := range testHeadEngineNames {
		engine := spec.Engine{Engine: engineName}
		summary := Summary{Engines: []spec.Engine{engine}}
		assert.False(t, summary.IsServerProcessEngine())
	}
}
