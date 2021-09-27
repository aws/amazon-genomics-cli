package spec

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/stretchr/testify/assert"
)

func TestProjectValidation_Valid(t *testing.T) {
	tests := map[string]struct {
		yaml string
	}{
		"valid": {
			yaml: `---
name: Demo
schemaVersion: 1
# I am a comment
workflows:
  hello:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/hello.wdl
  read:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/read.wdl
  haplotype:
    type:
      language: wdl
      # another comment
      version: 1.0
    sourceURL: workflows/haplotypecaller-gvcf-gatk4.wdl
  words-with-vowels:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/words-with-vowels.wdl
data:
  - location: s3://gatk-test-data
    readOnly: true
    # cool data
  - location: s3://broad-references
    readOnly: true
contexts:
    default:
        engines:
            - type: wdl
              engine: cromwell
    frugal:
        requestSpotInstances: true
        instanceTypes:
            - t3
        engines:
            - type: wdl
              engine: cromwell
    dev:
        requestSpotInstances: true
        engines:
            - type: wdl
              engine: cromwell
    prod:
        requestSpotInstances: false
        instanceTypes:
            - c6gd.16xlarge
        engines:
            - type: wdl
              engine: cromwell
`,
		},
		"defaultContext": {
			yaml: `---
name: foo
schemaVersion: 1
contexts:
    myContext:
        engines:
            - type: wdl
              engine: cromwell`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateProject([]byte(tt.yaml))
			assert.NoError(t, err, "Expected valid project spec")
		})
	}
}

func TestProjectValidation_Invalid(t *testing.T) {
	tests := map[string]struct {
		yaml       string
		errMessage string
	}{
		"two engines in a context": {
			yaml: `---
name: Demo
schemaVersion: 1
workflows:
  hello:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/hello.wdl
contexts:
    twoEngines: 
        engines:
            - type: wdl
              engine: cromwell
            - type: nextflow
              engine: nextflow
`,
			errMessage: "\n\t1. contexts.twoEngines.engines: Array must have at most 1 items\n",
		},
		"missing engine": {
			yaml: `---
name: Demo
schemaVersion: 1
workflows:
  hello:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/hello.wdl
contexts:
    noEngines: {}
`,
			errMessage: "\n\t1. contexts.noEngines: engines is required\n",
		},
		"empty project": {
			yaml:       ``,
			errMessage: "\n\t1. (root): Invalid type. Expected: object, given: null\n",
		},
		"no context": {
			yaml:       `name: foo`,
			errMessage: "\n\t1. (root): contexts is required\n",
		},
		"missingWorkflowFields": {
			yaml: `---
name: foo
schemaVersion: 1
contexts: 
    default:
        engines:
            - type: wdl
              engine: cromwell
workflows:
    foo: {}`,
			errMessage: "\n\t1. workflows.foo: type is required\n\t2. workflows.foo: sourceURL is required\n",
		},
		"missingDataFields": {
			yaml: `---
name: foo
schemaVersion: 1
contexts: 
    default:
        engines:
            - type: wdl
              engine: cromwell
data:
  - readOnly: true`,
			errMessage: "\n\t1. data.0: location is required\n",
		},
		"badSchemaVersion": {
			yaml: `---
name: foo
schemaVersion: 0
contexts:
    myContext:
        engines:
            - type: wdl
              engine: cromwell`,
			errMessage: "\n\t1. schemaVersion: Must be greater than or equal to 1\n",
		},
		"missingEngineNameInContext": {
			yaml: `---
name: Demo
schemaVersion: 1
contexts:
    default:
        engines:
            - type: wdl
`,
			errMessage: "\n\t1. contexts.default.engines.0: engine is required\n",
		},
		"missingEngineTypeInContext": {
			yaml: `---
name: Demo
schemaVersion: 1
contexts:
    default:
        engines:
            - engine: cromwell
`,
			errMessage: "\n\t1. contexts.default.engines.0: type is required\n",
		},
		"zeroLengthTypeInContext": {
			yaml: `---
name: Demo
schemaVersion: 1
contexts:
    default:
        engines:
            - type: ''
              engine: 'cromwell'
`,
			errMessage: "\n\t1. contexts.default.engines.0.type: String length must be greater than or equal to 1\n",
		},
		"zeroLengthEngineNameInContext": {
			yaml: `---
name: Demo
schemaVersion: 1
contexts:
    default:
        engines:
            - type: 'wdl'
              engine: ''
`,
			errMessage: "\n\t1. contexts.default.engines.0.engine: String length must be greater than or equal to 1\n",
		},
		"invalidProjectName": {
			yaml: `---
name: MyEmoji_😬_Project!
schemaVersion: 1
contexts:
    default:
        engines:
            - type: wdl
              engine: cromwell
`,
			errMessage: "\n\t1. name: Does not match pattern '^[A-Za-z0-9]+$'\n",
		},
		"invalidContextName": {
			yaml: `---
name: Demo
schemaVersion: 1
contexts:
    not-default:
        engines:
            - type: wdl
              engine: cromwell
`,
			errMessage: "\n\t1. contexts: Additional property not-default is not allowed\n",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateProject([]byte(tt.yaml))
			if assert.Error(t, err, "Expected invalid project spec") {
				var aErr clierror.ActionableError
				if errors.As(err, &aErr) {
					assert.EqualError(t, aErr.Cause, tt.errMessage)
					return
				}
				assert.Failf(t, "Unexpected error type", "Expected an instance of clierror.ActionableError, got: %#v", err)
			}
		})
	}
}
