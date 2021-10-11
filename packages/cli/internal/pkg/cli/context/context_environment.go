package context

import "strconv"

const (
	contextDir = "context"
)

type contextEnvironment struct {
	ProjectName      string
	ContextName      string
	UserId           string
	UserEmail        string
	OutputBucketName string

	EngineName            string
	EngineDesignation     string
	EngineRepository      string
	EngineHealthCheckPath string

	AdapterName        string
	AdapterDesignation string
	AdapterRepository  string

	ArtifactBucketName   string
	ReadBucketArns       string
	ReadWriteBucketArns  string
	InstanceTypes        string
	ResourceType         string
	MaxVCpus             int
	RequestSpotInstances bool
}

func (input contextEnvironment) ToEnvironmentList() []string {
	return environmentMapToList(map[string]string{
		"PROJECT":       input.ProjectName,
		"CONTEXT":       input.ContextName,
		"USER_ID":       input.UserId,
		"USER_EMAIL":    input.UserEmail,
		"OUTPUT_BUCKET": input.OutputBucketName,

		"ENGINE_NAME":              input.EngineName,
		"ENGINE_DESIGNATION":       input.EngineDesignation,
		"ENGINE_REPOSITORY":        input.EngineRepository,
		"ENGINE_HEALTH_CHECK_PATH": input.EngineHealthCheckPath,

		"ADAPTER_NAME":        input.AdapterName,
		"ADAPTER_DESIGNATION": input.AdapterDesignation,
		"ADAPTER_REPOSITORY":  input.AdapterRepository,

		"ARTIFACT_BUCKET":              input.ArtifactBucketName,
		"READ_BUCKET_ARNS":             input.ReadBucketArns,
		"READ_WRITE_BUCKET_ARNS":       input.ReadWriteBucketArns,
		"BATCH_COMPUTE_INSTANCE_TYPES": input.InstanceTypes,
		"MAX_V_CPUS":                   strconv.Itoa(input.MaxVCpus),
		"REQUEST_SPOT_INSTANCES":       strconv.FormatBool(input.RequestSpotInstances),
	})
}
