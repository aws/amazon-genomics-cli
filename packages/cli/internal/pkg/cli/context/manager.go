package context

import (
	"net/url"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ssm"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/rs/zerolog/log"
)

const (
	listDelimiter     = ","
	artifactParameter = "installed-artifacts/s3-root-url"
	cdkAppsDirBase    = ".agc/cdk/apps"
)

//nolint:structcheck
type baseProps struct {
	projectSpec spec.Project
	contextSpec spec.Context
	userId      string
	userEmail   string
	homeDir     string
}

//nolint:structcheck
type contextProps struct {
	readBuckets      []string
	readWriteBuckets []string
	outputBucket     string
	artifactBucket   string
	artifactUrl      string
	contextEnv       contextEnvironment
}

//nolint:structcheck
type infoProps struct {
	contextStackInfo cfn.StackInfo
	contextStatus    Status
}

//nolint:structcheck
type listProps struct {
	contexts map[string]Summary
}

type Manager struct {
	Cdk     cdk.Interface
	Cfn     cfn.Interface
	Project storage.ProjectClient
	Config  storage.ConfigClient
	Ssm     ssm.Interface

	baseProps
	contextProps
	infoProps
	listProps
	err error
}

func NewManager(profile string) *Manager {
	projectClient, err := storage.NewProjectClient()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create Project client for context manager")
	}
	homeDir, err := config.DetermineHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to determine home directory")
	}
	configClient, err := config.NewConfigClient()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create config client for context manager")
	}
	return &Manager{
		Cdk:       aws.CdkClient(profile),
		Cfn:       aws.CfnClient(profile),
		Project:   projectClient,
		Config:    configClient,
		Ssm:       aws.SsmClient(profile),
		baseProps: baseProps{homeDir: homeDir},
	}
}

func (m *Manager) readProjectSpec() {
	if m.err != nil {
		return
	}
	m.projectSpec, m.err = m.Project.Read()
}

func (m *Manager) readContextSpec(contextName string) {
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

func (m *Manager) setDataBuckets() {
	if m.err != nil {
		return
	}
	for _, dataItem := range m.projectSpec.Data {
		s3Arn, err := s3.UriToArn(dataItem.Location)
		if err != nil {
			m.err = err
			return
		}
		if dataItem.ReadOnly {
			m.readBuckets = append(m.readBuckets, s3Arn)
		} else {
			m.readWriteBuckets = append(m.readWriteBuckets, s3Arn)
		}
	}
}

func (m *Manager) setArtifactUrl() {
	if m.err != nil {
		return
	}
	m.artifactUrl, m.err = m.Ssm.GetCommonParameter(artifactParameter)
	if m.err != nil {
		m.err = actionableerror.FindSuggestionForError(m.err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
}

func (m *Manager) setArtifactBucket() {
	if m.err != nil {
		return
	}
	parsedUrl, err := url.Parse(m.artifactUrl)
	if err != nil {
		m.err = err
		return
	}
	m.artifactBucket = parsedUrl.Host
}

func (m *Manager) setOutputBucket() {
	if m.err != nil {
		return
	}
	m.outputBucket, m.err = m.Ssm.GetOutputBucket()
	if m.err != nil {
		m.err = actionableerror.FindSuggestionForError(m.err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
}

func (m *Manager) setTaskContext(contextName string) {
	if m.err != nil {
		return
	}

	m.contextEnv = contextEnvironment{
		ProjectName:          m.projectSpec.Name,
		ContextName:          contextName,
		UserId:               m.userId,
		UserEmail:            m.userEmail,
		OutputBucketName:     m.outputBucket,
		ArtifactBucketName:   m.artifactBucket,
		ReadBucketArns:       strings.Join(m.readBuckets, listDelimiter),
		ReadWriteBucketArns:  strings.Join(m.readWriteBuckets, listDelimiter),
		InstanceTypes:        strings.Join(m.contextSpec.InstanceTypes, listDelimiter),
		MaxVCpus:             m.contextSpec.MaxVCpus,
		RequestSpotInstances: m.contextSpec.RequestSpotInstances,
	}
}

func (m *Manager) setContextEnv(contextName string) {
	if m.err != nil {
		return
	}

	context, err := m.projectSpec.GetContext(contextName)
	if err != nil {
		m.err = err
		return
	}

	m.contextEnv = contextEnvironment{
		ProjectName:          m.projectSpec.Name,
		ContextName:          contextName,
		UserId:               m.userId,
		UserEmail:            m.userEmail,
		OutputBucketName:     m.outputBucket,
		ArtifactBucketName:   m.artifactBucket,
		ReadBucketArns:       strings.Join(m.readBuckets, listDelimiter),
		ReadWriteBucketArns:  strings.Join(m.readWriteBuckets, listDelimiter),
		InstanceTypes:        strings.Join(m.contextSpec.InstanceTypes, listDelimiter),
		MaxVCpus:             m.contextSpec.MaxVCpus,
		RequestSpotInstances: m.contextSpec.RequestSpotInstances,
		// TODO: we default to a single engine in a context for now
		// need to allow for multiple engines in the same context
		EngineName:        context.Engines[0].Engine,
		EngineDesignation: context.Engines[0].Engine,
	}
}

func (m *Manager) readConfig() {
	if m.err != nil {
		return
	}
	m.userId, m.err = m.Config.GetUserId()
	if m.err != nil {
		return
	}
	m.userEmail, m.err = m.Config.GetUserEmailAddress()
}
