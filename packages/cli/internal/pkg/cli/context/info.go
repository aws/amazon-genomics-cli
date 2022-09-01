package context

import (
	"reflect"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
)

var serverProcessEngines = map[string]bool{constants.CROMWELL: true}

type Summary struct {
	Name          string
	MaxVCpus      int
	IsSpot        bool
	InstanceTypes []string
	Engines       []spec.Engine
}

func (s Summary) IsEmpty() bool {
	return reflect.ValueOf(s).IsZero()
}

// IsServerProcessEngine Does the workflow engine run as a server process? A server process engine has one to many
// mapping with workflow runs. The engine can be used to run multiple workflows and the process is re-used and long running.
func (s *Summary) IsServerProcessEngine() bool {
	return serverProcessEngines[s.Engines[0].Engine]
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

func (d Detail) IsEmpty() bool {
	return reflect.ValueOf(d).IsZero()
}
