package spec

type Engine struct {
	Type   string `yaml:"type"`
	Engine string `yaml:"engine"`
}

type Context struct {
	InstanceTypes        []string `yaml:"instanceTypes,omitempty"`
	RequestSpotInstances bool     `yaml:"requestSpotInstances,omitempty"`
	Engines              []Engine `yaml:"engines"`
}
