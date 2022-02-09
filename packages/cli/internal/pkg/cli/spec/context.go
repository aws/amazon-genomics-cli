package spec

import "fmt"

const DefaultMaxVCpus = 256
const DefaultFSProvisionedThroughput = 0

type FSConfig struct {
	FSProvisionedThroughput int `yaml:"provisionedThroughput"`
}
type Filesystem struct {
	FSType        string   `yaml:"fsType"`
	Configuration FSConfig `yaml:"configuration,omitempty"`
}
type Engine struct {
	Type       string     `yaml:"type"`
	Engine     string     `yaml:"engine"`
	Filesystem Filesystem `yaml:"filesystem,omitempty"`
}
type Context struct {
	InstanceTypes        []string `yaml:"instanceTypes,omitempty"`
	RequestSpotInstances bool     `yaml:"requestSpotInstances,omitempty"`
	MaxVCpus             int      `yaml:"maxVCpus,omitempty"`
	Engines              []Engine `yaml:"engines"`
}

func (filesystem *Filesystem) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type defValFilesystem Filesystem
	defaults := defValFilesystem{Configuration: FSConfig{FSProvisionedThroughput: DefaultFSProvisionedThroughput}}
	if err := unmarshal(&defaults); err != nil {
		return err
	}

	*filesystem = Filesystem(defaults)
	switch filesystem.FSType {
	case "S3", "EFS", "":
		return nil
	default:
		return fmt.Errorf("filesystem %s is invalid. Options are `S3` and `EFS`", filesystem.FSType)
	}
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
