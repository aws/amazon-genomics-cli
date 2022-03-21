package context

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"reflect"
	"strings"
)

var serverEngines = map[string]bool{"cromwell": true}
var headEngines = map[string]bool{"nextflow": true, "miniwdl": true, "snakemake": true}

type Summary struct {
	Name          string
	MaxVCpus      int
	IsSpot        bool
	InstanceTypes []string
	Engines       []spec.Engine
}

func (i Summary) IsEmpty() bool {
	return reflect.ValueOf(i).IsZero()
}

//IsServerProcessEngine Does the workflow engine run as a server process? A server process has a potential one to many
// mapping with workflow runs.
func (i *Summary) IsServerProcessEngine() bool {
	return serverEngines[strings.ToLower(strings.TrimSpace(i.Engines[0].Engine))]
}

//IsHeadProcessEngine Does the workflow engine run as a head process? A head process has one to one
// mapping with workflow runs. Processes are not reused.
func (i *Summary) IsHeadProcessEngine() bool {
	return headEngines[strings.ToLower(strings.TrimSpace(i.Engines[0].Engine))]
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
