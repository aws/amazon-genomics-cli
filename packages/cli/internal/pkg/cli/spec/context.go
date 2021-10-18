package spec

const DefaultMaxVCpus = 256

type Engine struct {
	Type   string `yaml:"type"`
	Engine string `yaml:"engine"`
}

type Context struct {
	InstanceTypes        []string `yaml:"instanceTypes,omitempty"`
	RequestSpotInstances bool     `yaml:"requestSpotInstances,omitempty"`
	MaxVCpus             int      `yaml:"maxVCpus,omitempty"`
	Engines              []Engine `yaml:"engines"`
}

func (context *Context) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type defValContext Context
	defaults := defValContext{MaxVCpus: DefaultMaxVCpus}
	if err := unmarshal(&defaults); err != nil {
		return err
	}

	*context = Context(defaults)
	return nil
}
