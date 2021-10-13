package types

type Context struct {
	Name                 string
	Status               string
	StatusReason         string
	MaxVCpus             int
	RequestSpotInstances bool
	InstanceTypes        []InstanceType
	Output               OutputLocation
}

type ContextInstance struct {
	Id          string
	Name        string
	RunStatus   string
	ErrorStatus string
	StartTime   string
	RunTime     string
	Info        string
}

type ContextName struct {
	Name string
}

type OutputLocation struct {
	Url string
}

type InstanceType struct {
	Value string
}
