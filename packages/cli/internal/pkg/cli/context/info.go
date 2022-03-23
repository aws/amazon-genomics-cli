package context

import (
	"reflect"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
)

var headEngines = map[string]bool{constants.NEXTFLOW: true, constants.MINIWDL: true, constants.SNAKEMAKE: true}

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

//IsHeadProcessEngine Does the workflow engine run as a head process? A head process has one to one
// mapping with workflow runs. Processes are not reused.
func (s *Summary) IsHeadProcessEngine() bool {
	return headEngines[s.Engines[0].Engine]
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
