package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ssm"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/awsresources"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/zipfile"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
	"github.com/aws/amazon-genomics-cli/internal/pkg/osutils"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes/option"
	"github.com/rs/zerolog/log"
	"github.com/rsc/wes_client"
)

var (
	compressToTmp                 = zipfile.CompressToTmp
	workflowZip                   = "workflow.zip"
	manifestFilename              = storage.ManifestFileName
	removeFile                    = os.Remove
	removeAll                     = os.RemoveAll
	osStat                        = os.Stat
	createTempDir                 = os.MkdirTemp
	copyFileRecursivelyToLocation = osutils.CopyFileRecursivelyToLocation
	writeToTmp                    = func(namePattern, content string) (string, error) {
		f, err := os.CreateTemp("", namePattern)
		if err != nil {
			return "", err
		}
		defer f.Close()
		_, err = f.WriteString(content)
		if err != nil {
			return "", err
		}
		return f.Name(), nil
	}
)

//nolint:structcheck
type baseProps struct {
	projectSpec spec.Project
	contextSpec spec.Context
	userId      string
}

//nolint:structcheck
type s3Props struct {
	bucketName      string
	baseWorkflowKey string
}

//nolint:structcheck
type runProps struct {
	runId                string
	workflowSpec         spec.Workflow
	workflowEngine       string
	parsedSourceURL      *url.URL
	isLocal              bool
	path                 string
	packPath             string
	workflowUrl          string
	manifestPath         string
	tempPath             string
	inputsPath           string
	input                Input
	optionFileUrl        string
	options              map[string]string
	arguments            []string
	attachments          []string
	workflowParams       map[string]string
	workflowEngineParams map[string]string
}

//nolint:structcheck
type listProps struct {
	workflows map[string]Summary
}

//nolint:structcheck
type wesProps struct {
	contextStackInfo cfn.StackInfo
	wesUrl           string
}

//nolint:structcheck
type instanceProps struct {
	instances           []InstanceSummary
	filteredInstances   []InstanceSummary
	instancesPerContext map[string][]*InstanceSummary
}

//nolint:structcheck
type instanceStopProps struct {
	instanceToStop ddb.WorkflowInstance
}

//nolint:structcheck
type taskProps struct {
	runContextName string
	runLog         wes_client.RunLog
}

//nolint:structcheck
type workflowOutputProps struct {
	instanceSummary       InstanceSummary
	workflowRunLogOutputs map[string]interface{}
}

type ManifestProps struct {
	MainWorkFlowURL string   `json:"mainWorkFlowURL"`
	InputFileURLs   []string `json:"inputFileURLs"`
	EngineOptions   string   `json:"engineOptions"`
}

type Manager struct {
	Project     storage.ProjectClient
	Config      storage.ConfigClient
	S3          s3.Interface
	Ssm         ssm.Interface
	Cfn         cfn.Interface
	Ddb         ddb.Interface
	Storage     storage.StorageClient
	InputClient storage.InputClient
	WesFactory  func(url string) (wes.Interface, error)

	wes wes.Interface
	baseProps
	wesProps
	s3Props
	runProps
	listProps
	instanceProps
	instanceStopProps
	taskProps
	workflowOutputProps
	err error
}

func NewManager(profile string) *Manager {
	projectClient, err := storage.NewProjectClient()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create Project client for workflow manager")
	}
	storageClient, err := storage.NewStorageInstance()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create storage client for workflow manager")
	}
	configClient, err := config.NewConfigClient()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create config client for workflow manager")
	}
	s3Client := aws.S3Client(profile)
	return &Manager{
		Project:     projectClient,
		Config:      configClient,
		Ssm:         aws.SsmClient(profile),
		Cfn:         aws.CfnClient(profile),
		S3:          s3Client,
		Ddb:         aws.DdbClient(profile),
		Storage:     storageClient,
		InputClient: storage.NewInputClient(s3Client),
		WesFactory:  func(url string) (wes.Interface, error) { return wes.New(url, profile) },
	}
}

func (m *Manager) readProjectSpec() {
	if m.err != nil {
		return
	}
	log.Debug().Msgf("reading project specification")
	m.projectSpec, m.err = m.Project.Read()
}

func (m *Manager) validateContextIsDeployed(contextName string) {
	if m.err != nil {
		return
	}
	log.Debug().Msgf("checking deployment status of '%s' context", contextName)
	if !m.isContextDeployed(contextName) && m.err == nil {
		m.err = fmt.Errorf("context '%s' is not deployed", contextName)
	}
}

func (m *Manager) setWorkflowSpec(workflowName string) {
	if m.err != nil {
		return
	}
	log.Debug().Msgf("reading specification of '%s' workflow", workflowName)
	workflowSpec, ok := m.projectSpec.Workflows[workflowName]
	if !ok {
		m.err = fmt.Errorf("workflow '%s' is not defined in Project '%s' specification", workflowName, m.projectSpec.Name)
		return
	}
	log.Debug().Msgf("workflow type: '%s' version: '%s', workflow source url: '%s'", workflowSpec.Type.Language, workflowSpec.Type.Version, workflowSpec.SourceURL)
	m.workflowSpec = workflowSpec
}

func (m *Manager) parseWorkflowLocation() {
	if m.err != nil {
		return
	}
	m.parsedSourceURL, m.err = url.Parse(m.workflowSpec.SourceURL)
	if m.err == nil {
		log.Debug().Msgf("parsed workflow location as '%s'", m.parsedSourceURL.String())
	}
}

func (m *Manager) isUploadRequired() bool {
	if m.err != nil {
		return false
	}
	scheme := strings.ToLower(m.parsedSourceURL.Scheme)
	m.isLocal = scheme == "" || scheme == "file"
	log.Debug().Msgf("workflow location is local? '%t', upload is required? '%t'", m.isLocal, m.isLocal)
	inputIncluded := m.inputsPath != ""
	log.Debug().Msgf("does input file exist? '%t'", inputIncluded)
	return m.isLocal || inputIncluded
}

func (m *Manager) setWorkflowPath() {
	if m.err != nil {
		return
	}
	projectLocation := m.Project.GetLocation()
	workflowPath := m.parsedSourceURL.Path
	m.path = filepath.Join(projectLocation, workflowPath)
	log.Debug().Msgf("workflow path is '%s", m.path)
}

func (m *Manager) packWorkflowPath() {
	if m.err != nil {
		return
	}

	fileInfo, err := osStat(m.path)
	if err != nil {
		m.err = err
		return
	}

	if fileInfo.IsDir() {
		log.Debug().Msgf("workflow path '%s' is a directory, packing contents ...", m.tempPath)
		defer m.deleteTempDir()
		log.Debug().Msgf("updating file references and loading packed content to '%s/%s'", m.bucketName, m.baseWorkflowKey)
		err = m.InputClient.UpdateInputReferencesAndUploadToS3(m.path, m.tempPath, m.bucketName, m.baseWorkflowKey)
		if err != nil {
			log.Error().Err(err)
			m.err = err
			return
		}
	} else {
		m.tempPath = m.path
	}

	m.packPath, m.err = compressToTmp(m.tempPath)
}

func (m *Manager) setOutputBucket() {
	if m.err != nil {
		return
	}
	m.bucketName, m.err = m.Ssm.GetOutputBucket()
	log.Debug().Msgf("using output bucket '%s'", m.bucketName)
}

func (m *Manager) setBaseObjectKey(contextName, workflowName string) {
	if m.err != nil {
		return
	}
	m.baseWorkflowKey = awsresources.RenderBucketContextKey(m.projectSpec.Name, m.userId, contextName, "workflow", workflowName)
	log.Debug().Msgf("workflow upload base object key is '%s'", m.baseWorkflowKey)
}

func (m *Manager) calculateFinalLocation() {
	if m.err != nil {
		return
	}
	if m.isLocal {
		m.workflowUrl = fmt.Sprintf("s3://%s/%s/%s", m.bucketName, m.baseWorkflowKey, workflowZip)
	} else {
		m.workflowUrl = m.workflowSpec.SourceURL
	}
	log.Debug().Msgf("workflow artifacts at '%s' will be used to run the workflow", m.workflowUrl)
}

func (m *Manager) uploadWorkflowToS3() {
	if m.err != nil {
		return
	}
	objectKey := fmt.Sprintf("%s/%s", m.baseWorkflowKey, workflowZip)
	log.Debug().Msgf("updloading '%s' to 's3://%s/%s", m.packPath, m.bucketName, objectKey)
	m.err = m.S3.UploadFile(m.bucketName, objectKey, m.packPath)
	if m.err != nil {
		m.err = fmt.Errorf("unable to upload s3://%s/%s: %w", m.bucketName, objectKey, m.err)
	}
}

func (m *Manager) readInput(inputUrl string) {
	if m.err != nil || inputUrl == "" {
		return
	}
	log.Debug().Msgf("Input file override URL: %s", inputUrl)
	m.inputsPath = osutils.StripFileURLPrefix(inputUrl) // We actually support only local files
	bytes, err := m.Storage.ReadAsBytes(inputUrl)
	log.Debug().Msgf("content is:\n'%s'", string(bytes))
	if err != nil {
		m.err = err
		return
	}
	var input Input
	if err := json.Unmarshal(bytes, &input); err != nil {
		m.err = err
		return
	}
	m.input = input
}

func (m *Manager) copyInputToTemp() {
	if m.err != nil || m.inputsPath == "" {
		return
	}
	log.Debug().Msgf("Copying input file to temp: %s", m.tempPath)
	absInputsPath, err := filepath.Abs(m.inputsPath)
	bytes, err := m.Storage.ReadAsBytes(m.inputsPath)

	if err != nil {
		m.err = err
		return
	}
	dest := filepath.Join(m.tempPath, filepath.Base(absInputsPath))
	err = m.Storage.WriteFromBytes(dest, bytes)
	if err != nil {
		m.err = err
		return
	}
}

func (m *Manager) parseInputToArguments() {
	if m.err != nil || m.input == nil {
		return
	}
	arguments := m.input.String()
	log.Debug().Msgf("arguments are: '%s'", arguments)
	m.arguments = []string{arguments}
}

func (m *Manager) initializeTempDir() {
	if m.err != nil {
		return
	}
	var err error
	m.tempPath, err = createTempDir("", "workflow_*")
	log.Debug().Msgf("created temp directory at: '%s'", m.tempPath)
	if err != nil {
		m.err = err
		return
	}
	log.Debug().Msgf("recursively copying content of '%s' to '%s'", m.path, m.tempPath)
	err = copyFileRecursivelyToLocation(m.tempPath, m.path)
	if err != nil {
		log.Error().Err(err)
		m.err = err
		return
	}
}

func (m *Manager) deleteTempDir() {
	if m.tempPath == "" {
		return
	}
	err := removeAll(m.tempPath)
	if err != nil {
		log.Warn().Msgf("Failed to delete temporary folder '%s'", m.tempPath)
	}
}

// writeTempManifest writes the inputsFile included in the command line to the temporary MANIFEST.json located in temp directory
// This function is only called if there is a path included in the command line with the --inputsFile flag
func (m *Manager) writeTempManifest() {
	if m.err != nil || m.inputsPath == "" {
		return
	}
	m.manifestPath = filepath.Join(m.tempPath, manifestFilename)
	log.Debug().Msgf("Reading temp manifest %s", m.manifestPath)
	bytes, err := m.Storage.ReadAsBytes(m.manifestPath)
	if err != nil {
		m.err = err
		return
	}
	var data ManifestProps
	if err := json.Unmarshal(bytes, &data); err != nil {
		m.err = err
		return
	}
	data.InputFileURLs = append(data.InputFileURLs, filepath.Base(m.inputsPath))
	bytes, err = json.Marshal(data)
	if err != nil {
		m.err = err
		return
	}
	err = m.Storage.WriteFromBytes(m.manifestPath, bytes)
	if err != nil {
		m.err = err
		return
	}
}

func (m *Manager) uploadInputsToS3() {
	if m.err != nil || m.input == nil {
		return
	}
	objectKey := awsresources.RenderBucketDataKey(m.projectSpec.Name, m.userId)
	absInputsPath, err := filepath.Abs(m.inputsPath)
	if err != nil {
		m.err = err
		return
	}
	baseLocation := filepath.Dir(absInputsPath)
	log.Debug().Msgf("moving local inputs from '%s' to s3://%s/%s and replacing paths with S3 paths", baseLocation, m.bucketName, objectKey)
	inputsWithS3Paths, err := m.InputClient.UpdateInputs(baseLocation, m.input, m.bucketName, objectKey)
	if err != nil {
		m.err = fmt.Errorf("unable to sync s3://%s/%s: %w", m.bucketName, objectKey, err)
		return
	}
	m.input = inputsWithS3Paths
}

func (m *Manager) readOptionFile(optionFileUrl string) {
	if m.err != nil || optionFileUrl == "" {
		return
	}
	log.Debug().Msgf("Option file override URL: %s", optionFileUrl)
	m.optionFileUrl = optionFileUrl
	bytes, err := m.Storage.ReadAsBytes(optionFileUrl)
	log.Debug().Msgf("with content:\n%s", string(bytes))
	if err != nil {
		m.err = err
		return
	}
	var options map[string]string
	if err := json.Unmarshal(bytes, &options); err != nil {
		m.err = err
		return
	}
	m.options = options
}

func (m *Manager) readConfig() {
	if m.err != nil {
		return
	}
	m.userId, m.err = m.Config.GetUserId()
	log.Debug().Msgf("current user id: '%s'", m.userId)
}

func (m *Manager) isContextDeployed(contextName string) bool {
	if m.err != nil {
		return false
	}
	engineStackName := awsresources.RenderContextStackName(m.projectSpec.Name, contextName, m.userId)
	status, err := m.Cfn.GetStackStatus(engineStackName)
	if err != nil {
		if errors.Is(err, cfn.StackDoesNotExistError) {
			return false
		}
		m.err = err
		return false
	}

	ok, activeStatusFlag := cfn.QueryableStacksMap[status]
	return ok && activeStatusFlag
}

func (m *Manager) setContext(contextName string) {
	if m.err != nil {
		return
	}

	log.Debug().Msgf("obtaining spec for context '%s'", contextName)
	contextSpec, err := m.projectSpec.GetContext(contextName)
	if err != nil {
		m.err = err
		return
	}
	m.contextSpec = contextSpec
}

func (m *Manager) setEngineForWorkflowType(contextName string) {
	if m.err != nil {
		return
	}
	enginesLen := len(m.contextSpec.Engines)
	if enginesLen == 0 {
		m.err = fmt.Errorf("context '%s' doesn't have any engines defined", contextName)
		return
	}
	if enginesLen > 1 {
		m.err = fmt.Errorf("only one engine per context is supported. Context '%s' has %d engines defined", contextName, enginesLen)
		return
	}
	log.Debug().Msgf("using engine '%s' from context '%s'", m.contextSpec.Engines[0].Engine, contextName)
	m.workflowEngine = m.contextSpec.Engines[0].Engine
}

func (m *Manager) setContextStackInfo(contextName string) {
	if m.err != nil {
		return
	}
	contextStackName := awsresources.RenderContextStackName(m.projectSpec.Name, contextName, m.userId)
	log.Debug().Msgf("using context infrastructure from cloudformation stack '%s'", contextStackName)
	m.contextStackInfo, m.err = m.Cfn.GetStackInfo(contextStackName)
}

func (m *Manager) setWesUrl() {
	if m.err != nil {
		return
	}
	wesUrl, ok := m.contextStackInfo.Outputs["WesUrl"]
	if !ok {
		m.err = fmt.Errorf("wes endpoint for workflow type '%s' is missing in engine stack '%s'",
			m.workflowSpec.Type.Language, m.contextStackInfo.Id)
		return
	}
	log.Debug().Msgf("workflow will be submitted to wes endpoint at '%s'", wesUrl)
	m.wesUrl = wesUrl
}

func (m *Manager) setWorkflowParameters() {
	if m.err != nil {
		return
	}
	m.workflowParams = make(map[string]string)
	if m.inputsPath == "" {
		return
	}
	m.workflowParams["workflowInputs"] = filepath.Base(m.attachments[0])
	log.Debug().Msgf("workflow parameter of 'workflowInputs' is '%s'", m.workflowParams["workflowInputs"])
}

func (m *Manager) setWorkflowEngineParameters() {
	if m.err != nil {
		return
	}
	m.workflowEngineParams = make(map[string]string)
	if m.optionFileUrl == "" {
		return
	}
	if m.options != nil {
		if m.workflowEngine == constants.NEXTFLOW || m.workflowEngine == constants.MINIWDL || m.workflowEngine == constants.SNAKEMAKE {
			m.err = fmt.Errorf("optionFile flag cannot be used with head node engines")
			return
		}
		m.workflowEngineParams = m.options
		log.Debug().Msgf("workflow engine parameters: %s", fmt.Sprint(m.workflowEngineParams))
	}
}

func (m *Manager) setWesClient() {
	if m.err != nil {
		return
	}
	log.Debug().Msgf("constructing API client for WES endpoint at '%s'", m.wesUrl)
	m.wes, m.err = m.WesFactory(m.wesUrl)
	if m.err != nil {
		m.err = fmt.Errorf("unable to configure client for WES endpoint: %w", m.err)
	}
}

func (m *Manager) saveAttachments() {
	if m.err != nil {
		return
	}

	namePattern := fmt.Sprintf("%s_*", filepath.Base(m.inputsPath))
	for _, arg := range m.arguments {
		fileName, err := writeToTmp(namePattern, arg)
		log.Debug().Msgf("saved attachment for argument '%s' to '%s'", arg, fileName)
		if err != nil {
			m.err = err
			return
		}
		m.attachments = append(m.attachments, fileName)
	}
}

func (m *Manager) cleanUpAttachments() {
	for _, attachment := range m.attachments {
		log.Debug().Msgf("cleaning up '%s'", attachment)
		err := removeFile(attachment)
		if err != nil {
			log.Warn().Msgf("Failed to clean up temporary file '%s': %s", attachment, err)
		}
	}
}

func (m *Manager) runWorkflow() {
	if m.err != nil {
		return
	}
	log.Debug().Msgf("calling WES client run_workflow with workflow_url: '%s', workflow_type: '%s', workflow_type_version: '%s', workflow_attachment: '[%s]', workflow_params: '%s', workflow_engine_parameters: '%s'",
		m.workflowUrl,
		m.workflowSpec.Type.Language,
		m.workflowSpec.Type.Version,
		strings.Join(m.attachments, ", "),
		fmt.Sprint(m.workflowParams),
		fmt.Sprint(m.workflowEngineParams))

	m.runId, m.err = m.wes.RunWorkflow(
		context.Background(),
		option.WorkflowUrl(m.workflowUrl),
		option.WorkflowType(m.workflowSpec.Type.Language),
		option.WorkflowTypeVersion(m.workflowSpec.Type.Version),
		option.WorkflowAttachment(m.attachments),
		option.WorkflowParams(m.workflowParams),
		option.WorkflowEngineParams(m.workflowEngineParams))
}

func (m *Manager) recordWorkflowRun(workflowName, contextName string) {
	if m.err != nil {
		return
	}
	log.Debug().Msgf("recording workflow run metadata for workflow run id '%s' to DynamodDB", m.runId)
	err := m.Ddb.WriteWorkflowInstance(context.Background(), ddb.WorkflowInstance{
		RunId:        m.runId,
		WorkflowName: workflowName,
		ContextName:  contextName,
		ProjectName:  m.projectSpec.Name,
		UserId:       m.userId,
	})
	if err != nil {
		log.Warn().Msgf("recording of run id failed: %s", err)
	}
}

func (m *Manager) cleanUpWorkflow() {
	if m.packPath != "" {
		log.Debug().Msgf("cleaning up '%s'", m.packPath)
		err := removeFile(m.packPath)
		if err != nil {
			log.Warn().Msgf("Failed to delete temporary file '%s'", m.packPath)
		}
	}
}

func (m *Manager) initWorkflows() {
	if m.err != nil {
		return
	}
	m.workflows = make(map[string]Summary)
}

func (m *Manager) readWorkflowsFromSpec() {
	if m.err != nil {
		return
	}
	log.Debug().Msgf("There are %d workflows specified in project '%s'", len(m.projectSpec.Workflows), m.projectSpec.Name)
	for name := range m.projectSpec.Workflows {
		log.Debug().Msgf("Workflow '%s' is defined in project '%s'", name, m.projectSpec.Name)
		m.workflows[name] = Summary{Name: name}
	}
}

func (m *Manager) stopWorkflowInstance(runId string) {
	if m.err != nil {
		return
	}

	err := m.wes.StopWorkflow(context.Background(), runId)
	if err != nil {
		log.Warn().Msgf("Unable to stop workflow instance '%s', is the workflow instance running? WES response '%s'", runId, err)
		m.err = err
		return
	}
}

func (m *Manager) loadInstancesByWorkflow(workflowName string, numInstances int) {
	if m.err != nil {
		return
	}
	instances, err := m.Ddb.ListWorkflowInstancesByName(context.Background(), m.projectSpec.Name, m.userId, workflowName, numInstances)
	if err != nil {
		m.err = err
		return
	}
	m.populateInstancesAndMapToContexts(instances)
}

func (m *Manager) loadInstancesByContext(contextName string, numInstances int) {
	if m.err != nil {
		return
	}
	instances, err := m.Ddb.ListWorkflowInstancesByContext(context.Background(), m.projectSpec.Name, m.userId, contextName, numInstances)
	if err != nil {
		m.err = err
		return
	}
	m.populateInstancesAndMapToContexts(instances)
}

func (m *Manager) loadInstances(numInstances int) {
	if m.err != nil {
		return
	}
	instances, err := m.Ddb.ListWorkflowInstances(context.Background(), m.projectSpec.Name, m.userId, numInstances)
	if err != nil {
		m.err = err
		return
	}
	m.populateInstancesAndMapToContexts(instances)
}

func (m *Manager) populateInstancesAndMapToContexts(workflowInstances []ddb.WorkflowInstance) {
	m.instancesPerContext = make(map[string][]*InstanceSummary)
	m.instances = make([]InstanceSummary, len(workflowInstances))
	for i, instance := range workflowInstances {
		key := instance.ContextName
		instanceSummary := InstanceSummary{
			Id:           instance.RunId,
			WorkflowName: instance.WorkflowName,
			ContextName:  instance.ContextName,
			SubmitTime:   instance.CreatedTime,
			Request:      instance.Request,
		}
		m.instances[i] = instanceSummary
		m.instancesPerContext[key] = append(m.instancesPerContext[key], &m.instances[i])
	}
}

func (m *Manager) loadInstance(instanceId string) {
	if m.err != nil {
		return
	}
	instance, err := m.Ddb.GetWorkflowInstanceById(context.Background(), m.projectSpec.Name, m.userId, instanceId)
	if err != nil {
		m.err = err
		return
	}
	m.populateInstancesAndMapToContexts([]ddb.WorkflowInstance{instance})
}

func (m *Manager) updateInstanceState(instance *InstanceSummary) {
	if m.err != nil || instance == nil {
		return
	}
	instance.State, m.err = m.wes.GetRunStatus(context.Background(), instance.Id)
}

func (m *Manager) setFilteredInstances() {
	if m.err != nil {
		return
	}
	for _, instance := range m.instances {
		if instance.State == "UNKNOWN" || instance.State == "" {
			log.Debug().Msgf("Workflow instance '%s' status is '%s', skipping", instance.Id, instance.State)
			continue
		}
		m.filteredInstances = append(m.filteredInstances, instance)
	}
}

func (m *Manager) renderWorkflowDetails(workflowName string) (Details, error) {
	if m.err != nil {
		return Details{}, m.err
	}
	details := Details{
		Name:         workflowName,
		TypeLanguage: m.workflowSpec.Type.Language,
		TypeVersion:  m.workflowSpec.Type.Version,
		Source:       m.workflowSpec.SourceURL,
	}
	return details, nil
}

func (m *Manager) updateInProject(instance *InstanceSummary) {
	if m.err != nil || instance == nil {
		return
	}
	_, instance.InProject = m.projectSpec.Workflows[instance.WorkflowName]
}

func (m *Manager) setRequest(instance *InstanceSummary) {
	if m.err != nil || instance == nil {
		return
	}
	if instance.Request != "" {
		log.Debug().Msgf("Instance request field is already set.")
		return
	}
	testRunLog, err := m.wes.GetRunLog(context.Background(), instance.Id)
	if err != nil {
		return
	}
	testReq := testRunLog.Request
	workflowEngineParamsJsonBytes, err := json.Marshal(testReq)
	if err != nil {
		return
	}
	instance.Request = string(workflowEngineParamsJsonBytes)
}

func (m *Manager) setInstanceToStop(runId string) {
	if m.err != nil {
		return
	}
	workflowInstance, err := m.Ddb.GetWorkflowInstanceById(context.Background(), m.projectSpec.Name, m.userId, runId)
	if err != nil {
		log.Warn().Msgf("No active workflow instance found for id '%s'", runId)
		m.err = err
		return
	}

	if workflowInstance.ContextName == "" {
		log.Error().Msgf("No context can be found for workflow instance '%s', check you have the correct workflow instance id", runId)
		m.err = fmt.Errorf("no context can be found for workflow instance '%s'", runId)
	}

	m.instanceToStop = workflowInstance
}

func (m *Manager) setInstanceSummary() {
	if m.err != nil {
		return
	}

	if len(m.instances) == 0 {
		m.err = fmt.Errorf("no instance summary found, check the workflow run id and check the workflow was run from the current project")
		log.Error().Err(m.err).Send()
		return
	}
	m.instanceSummary = m.instances[0]
}

func (m *Manager) setWorkflowRunLogOutputs() {
	if m.err != nil {
		return
	}
	runLog, err := m.wes.GetRunLog(context.Background(), m.instanceSummary.Id)
	if err != nil {
		m.err = err
		return
	}
	m.workflowRunLogOutputs = runLog.Outputs
}
