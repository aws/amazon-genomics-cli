package context

import "reflect"

type Summary struct {
	Name          string
	MaxVCpus      int
	IsSpot        bool
	InstanceTypes []string
}

func (i Summary) IsEmpty() bool {
	return reflect.ValueOf(i).IsZero()
}

type Detail struct {
	Summary
	Status             Status
	StatusReason       string
	BucketLocation     string
	WesUrl             string
	WesLogGroupName    string
	EngineLogGroupName string
	AccessLogGroupName string
}

type Instance struct {
	ContextName            string
	ContextStatus          Status
	ContextReason          string
	IsDefinedInProjectFile bool
}

func (i Detail) IsEmpty() bool {
	return reflect.ValueOf(i).IsZero()
}
