package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestProjectYaml(t *testing.T) {
	tests := map[string]struct {
		obj  Project
		yaml string
	}{
		"example": {
			obj: Project{
				Name:          "Demo",
				SchemaVersion: LatestVersion,
				Workflows: map[string]Workflow{
					"haplotype": {Type: WorkflowType{"wdl", "1.42"}, SourceURL: "file://workflows/haplotypecaller-gvcf-gatk4.wdl"},
					"hello":     {Type: WorkflowType{"wdl", "1.5"}, SourceURL: "file://workflows/hello.wdl"},
					"read":      {Type: WorkflowType{"wdl", "a.b"}, SourceURL: "file://workflows/read.wdl"},
				},
				Data: []Data{
					{Location: "s3://mybucket", ReadOnly: true},
				},
				Contexts: map[string]Context{
					"testContext": {
						MaxVCpus: 256,
						Engines: []Engine{
							{Type: "wdl", Engine: "cromwell"},
						},
					},
				},
			},
			yaml: `name: Demo
schemaVersion: 1
workflows:
    haplotype:
        type:
            language: wdl
            version: "1.42"
        sourceURL: file://workflows/haplotypecaller-gvcf-gatk4.wdl
    hello:
        type:
            language: wdl
            version: "1.5"
        sourceURL: file://workflows/hello.wdl
    read:
        type:
            language: wdl
            version: a.b
        sourceURL: file://workflows/read.wdl
data:
    - location: s3://mybucket
      readOnly: true
contexts:
    testContext:
        maxVCpus: 256
        engines:
            - type: wdl
              engine: cromwell
`,
		},
		"empty": {
			obj: Project{},
			yaml: `name: ""
schemaVersion: 0
`,
		},
		"complex": {
			obj: Project{
				Name:          "Complex",
				SchemaVersion: LatestVersion,
				Workflows: map[string]Workflow{
					"wf1": {Type: WorkflowType{"wdl", "1.0"}, SourceURL: "file://workflows/wf1.wdl"},
					"wf2": {Type: WorkflowType{"wdl", "abc.xyz"}, SourceURL: "s3://my-wf-bucket/wf2.wdl"},
				},
				Data: []Data{
					{Location: "s3://mybucket/a", ReadOnly: true},
					{Location: "s3://mybucket/b", ReadOnly: true},
					{Location: "s3://myotherbucket"},
				},
				Contexts: map[string]Context{
					"ctx1": {
						MaxVCpus: 256,
						Engines: []Engine{
							{Type: "wdl", Engine: "miniwdl"},
						},
					},
					"ctx2": {
						MaxVCpus: 256,
						Engines: []Engine{
							{Type: "nextflow", Engine: "nextflow"},
						},
					},
				},
			},
			yaml: `name: Complex
schemaVersion: 1
workflows:
    wf1:
        type:
            language: wdl
            version: "1.0"
        sourceURL: file://workflows/wf1.wdl
    wf2:
        type:
            language: wdl
            version: abc.xyz
        sourceURL: s3://my-wf-bucket/wf2.wdl
data:
    - location: s3://mybucket/a
      readOnly: true
    - location: s3://mybucket/b
      readOnly: true
    - location: s3://myotherbucket
contexts:
    ctx1:
        maxVCpus: 256
        engines:
            - type: wdl
              engine: miniwdl
    ctx2:
        maxVCpus: 256
        engines:
            - type: nextflow
              engine: nextflow
`,
		},
	}

	for name, tt := range tests {
		t.Run(name+"Unmarshal", func(t *testing.T) {
			expected := tt.obj
			actual := Project{}
			err := yaml.Unmarshal([]byte(tt.yaml), &actual)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
		t.Run(name+"Marshal", func(t *testing.T) {
			expected := tt.yaml
			bytes, err := yaml.Marshal(tt.obj)
			require.NoError(t, err)
			actual := string(bytes)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestProjectDefaults(t *testing.T) {
	const yamlStr = `
name: DefaultTest
schemaVersion: 1
contexts:
    context:
        engines:
            - type: wdl
              engine: cromwell
`

	t.Run("ContextDefaults", func(t *testing.T) {
		result := Project{}
		err := yaml.Unmarshal([]byte(yamlStr), &result)
		require.NoError(t, err)
		assert.Equal(t, result.Contexts["context"].MaxVCpus, DefaultMaxVCpus)
	})
}

func TestGetContext(t *testing.T) {
	type args struct {
		projectSpec Project
		contextName string
	}
	tests := []struct {
		name                 string
		args                 args
		expectedContext      Context
		expectedErrorMessage string
	}{
		{
			name: "Unknown context name",
			args: args{
				projectSpec: Project{
					Name: "myProject",
					Contexts: map[string]Context{
						"ctx1": {
							Engines: []Engine{
								{Type: "wdl", Engine: "miniwdl"},
							},
						},
						"ctx2": {
							Engines: []Engine{
								{Type: "nextflow", Engine: "nextflow"},
							},
						},
					},
				},
				contextName: "badContextName",
			},
			expectedErrorMessage: "context 'badContextName' is not defined in Project 'myProject' specification",
		},
		{
			name: "Existing context name ",
			args: args{
				projectSpec: Project{
					Name: "Complex",
					Contexts: map[string]Context{
						"ctx1": {
							Engines: []Engine{
								{Type: "wdl", Engine: "miniwdl"},
							},
						},
						"ctx2": {
							Engines: []Engine{
								{Type: "nextflow", Engine: "nextflow"},
							},
						},
					},
				},
				contextName: "ctx1",
			},
			expectedContext: Context{
				Engines: []Engine{
					{Type: "wdl", Engine: "miniwdl"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, err := tt.args.projectSpec.GetContext(tt.args.contextName)
			if tt.expectedErrorMessage != "" {
				assert.Error(t, err, tt.expectedErrorMessage)
			} else {
				assert.Equal(t, tt.expectedContext, context)
			}
		})
	}
}
