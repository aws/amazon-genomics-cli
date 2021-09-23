package spec

type Data struct {
	Location string `yaml:"location"`
	ReadOnly bool   `yaml:"readOnly,omitempty"`
}
