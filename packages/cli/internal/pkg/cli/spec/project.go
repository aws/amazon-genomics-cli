package spec

const LatestVersion = 1

type Project struct {
	Name          string              `yaml:"name"`
	SchemaVersion int                 `yaml:"schemaVersion"`
	Workflows     map[string]Workflow `yaml:"workflows,omitempty"`
	Data          []Data              `yaml:"data,omitempty"`
	Contexts      map[string]Context  `yaml:"contexts,omitempty"`
}
