package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes/option"
	"github.com/rs/zerolog/log"
	"github.com/rsc/wes_client"
)

var (
	compressToTmp = zipfile.CompressToTmp
	removeFile    = os.Remove
	writeToTmp    = func(namePattern, content string) (string, error) {
		f, err := ioutil.TempFile("", namePattern)
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
	chdir = os.Chdir
)

//nolint:structcheck
type baseProps struct {
	projectSpec spec.Project
	contextSpec spec.Context
	userId      string
}

//nolint:structcheck
type s3Props struct {
	bucketName string
	objectKey  string
}

//nolint:structcheck
type runProps struct {
	runId           string
	workflowSpec    spec.Workflow
	workflowEngine  string
	parsedSourceURL *url.URL
	isLocal         bool
	packPath        string
	workflowUrl     string
	inputUrl        string
	input           Input
	arguments       []string
	attachments     []string
	workflowParams  map[string]string
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

type Manager struct {
	Project    storage.ProjectClient
	Config     storage.ConfigClient
	S3         s3.Interface
	Ssm        ssm.Interface
	Cfn        cfn.Interface
	Ddb        ddb.Interface
	Storage    storage.StorageClient
	WesFactory func(url string) (wes.Interface, error)

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
	return &Manager{
		Project:    projectClient,
		Config:     configClient,
		Ssm:        aws.SsmClient(profile),
		Cfn:        aws.CfnClient(profile),
		S3:         aws.S3Client(profile),
		Ddb:        aws.DdbClient(profile),
		Storage:    storageClient,
		WesFactory: func(url string) (wes.Interface, error) { return wes.New(url, profile) },
	}
}

func (m *Manager) readProjectSpec() {
	if m.err != nil {
		return
	}
	m.projectSpec, m.err = m.Project.Read()
}

func (m *Manager) validateContextIsDeployed(contextName string) {
	if m.err != nil {
		return
	}
	if !m.isContextDeployed(contextName) && m.err == nil {
		m.err = fmt.Errorf("context '%s' is not deployed", contextName)
	}
}

func (m *Manager) chdirIntoProject() {
	if m.err != nil {
		return
	}
	m.err = chdir(m.Project.GetLocation())
}

func (m *Manager) setWorkflowSpec(workflowName string) {
	if m.err != nil {
		return
	}
	workflowSpec, ok := m.projectSpec.Workflows[workflowName]
	if !ok {
		m.err = fmt.Errorf("workflow '%s' is not defined in Project '%s' specification", workflowName, m.projectSpec.Name)
		return
	}
	m.workflowSpec = workflowSpec
}

func (m *Manager) parseWorkflowLocation() {
	if m.err != nil {
		return
	}
	m.parsedSourceURL, m.err = url.Parse(m.workflowSpec.SourceURL)
}

func (m *Manager) isUploadRequired() bool {
	if m.err != nil {
		return false
	}
	scheme := strings.ToLower(m.parsedSourceURL.Scheme)
	m.isLocal = scheme == "" || scheme == "file"
	return m.isLocal
}

func (m *Manager) packWorkflowFiles() {
	if m.err != nil {
		return
	}
	m.packPath, m.err = compressToTmp(m.parsedSourceURL.Path)
}

func (m *Manager) setOutputBucket() {
	if m.err != nil {
		return
	}
	m.bucketName, m.err = m.Ssm.GetOutputBucket()
}

func (m *Manager) setObjectKey(contextName, workflowName string) {
	if m.err != nil {
		return
	}
	m.objectKey = awsresources.RenderBucketContextKey(m.projectSpec.Name, m.userId, contextName, "workflow", workflowName, "workflow.zip")
}

func (m *Manager) calculateFinalLocation() {
	if m.err != nil {
		return
	}
	if m.isLocal {
		m.workflowUrl = fmt.Sprintf("s3://%s/%s", m.bucketName, m.objectKey)
	} else {
		m.workflowUrl = m.workflowSpec.SourceURL
	}
}

func (m *Manager) uploadWorkflowToS3() {
	if m.err != nil {
		return
	}
	m.err = m.S3.UploadFile(m.bucketName, m.objectKey, m.packPath)
}

func (m *Manager) readInput(inputUrl string) {
	if m.err != nil || inputUrl == "" {
		return
	}
	m.inputUrl = inputUrl
	bytes, err := m.Storage.ReadAsBytes(inputUrl)
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

func (m *Manager) parseInputToArguments() {
	if m.err != nil || m.input == nil {
		return
	}
	arguments, err := m.input.ToString()
	if err != nil {
		m.err = err
		return
	}
	m.arguments = []string{arguments}
}

func (m *Manager) uploadInputsToS3() {
	if m.err != nil || m.input == nil {
		return
	}
	m.input.MapInputUrls(m.uploadInputFileToS3)
}

func (m *Manager) uploadInputFileToS3(inputKey inputKey, fileUrl inputUrl) inputUrl {
	if m.err != nil {
		return fileUrl
	}
	parsedURL, err := url.Parse(fileUrl)
	if err != nil {
		m.err = err
		return fileUrl
	}
	scheme := strings.ToLower(parsedURL.Scheme)
	isLocal := scheme == "" || scheme == "file"
	if !isLocal {
		return fileUrl
	}
	objectKey := awsresources.RenderBucketDataKey(m.projectSpec.Name, m.userId, inputKey, filepath.Base(parsedURL.Path))
	absFileUrl, err := toAbsPath(filepath.Dir(m.inputUrl), fileUrl)
	if err != nil {
		m.err = err
		return fileUrl
	}
	err = m.S3.SyncFile(m.bucketName, objectKey, absFileUrl)
	if err != nil {
		m.err = err
		return fileUrl
	}
	return fmt.Sprintf("s3://%s/%s", m.bucketName, objectKey)
}

func toAbsPath(basePath, somePath string) (string, error) {
	if filepath.IsAbs(somePath) {
		return somePath, nil
	}
	return filepath.Abs(filepath.Join(basePath, somePath))
}

func (m *Manager) readConfig() {
	if m.err != nil {
		return
	}
	m.userId, m.err = m.Config.GetUserId()
}

func (m *Manager) isContextDeployed(contextName string) bool {
	if m.err != nil {
		return false
	}
	engineStackName := awsresources.RenderContextStackName(m.projectSpec.Name, contextName, m.userId)
	_, err := m.Cfn.GetStackStatus(engineStackName)
	if err != nil {
		if errors.Is(err, cfn.StackDoesNotExistError) {
			return false
		}
		m.err = err
		return false
	}
	return true
}

func (m *Manager) setContext(contextName string) {
	if m.err != nil {
		return
	}

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
	m.workflowEngine = m.contextSpec.Engines[0].Engine
}

func (m *Manager) setContextStackInfo(contextName string) {
	if m.err != nil {
		return
	}
	contextStackName := awsresources.RenderContextStackName(m.projectSpec.Name, contextName, m.userId)
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
	m.wesUrl = wesUrl
}

func (m *Manager) setWorkflowParameters() {
	if m.err != nil {
		return
	}
	m.workflowParams = make(map[string]string)
	if m.inputUrl == "" {
		return
	}
	m.workflowParams["workflowInputs"] = filepath.Base(m.attachments[0])
}

func (m *Manager) setWesClient() {
	if m.err != nil {
		return
	}
	m.wes, m.err = m.WesFactory(m.wesUrl)
}

func (m *Manager) saveAttachments() {
	if m.err != nil {
		return
	}

	for _, arg := range m.arguments {
		namePattern := fmt.Sprintf("%s_*", filepath.Base(m.inputUrl))
		fileName, err := writeToTmp(namePattern, arg)
		if err != nil {
			m.err = err
			return
		}
		m.attachments = append(m.attachments, fileName)
	}
}

func (m *Manager) cleanUpAttachments() {
	for _, attachment := range m.attachments {
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
	m.runId, m.err = m.wes.RunWorkflow(
		context.Background(),
		option.WorkflowUrl(m.workflowUrl),
		option.WorkflowType(m.workflowSpec.Type.Language),
		option.WorkflowTypeVersion(m.workflowSpec.Type.Version),
		option.WorkflowAttachment(m.attachments),
		option.WorkflowParams(m.workflowParams))
}

func (m *Manager) recordWorkflowRun(workflowName, contextName string) {
	if m.err != nil {
		return
	}
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
