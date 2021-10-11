package context

import (
	"errors"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/awsresources"
)

func (m *Manager) Info(contextName string) (Detail, error) {
	m.readProjectSpec()
	m.readConfig()
	m.setOutputBucket()
	m.setContextStackInfo(contextName)
	m.setContextEnv(contextName)
	m.parseContextStatus()
	return m.buildContextInfo(contextName)
}

func (m *Manager) setContextStackInfo(contextName string) {
	if m.err != nil {
		return
	}
	engineStackName := awsresources.RenderContextStackName(m.projectSpec.Name, contextName, m.userId)
	m.contextStackInfo, m.err = m.Cfn.GetStackInfo(engineStackName)
	if errors.Is(m.err, cfn.StackDoesNotExistError) {
		m.err = nil
	}
}

func (m *Manager) parseContextStatus() {
	if m.err != nil {
		return
	}
	status := mapStackToStatus(m.contextStackInfo.Status)
	switch {
	case status.IsFailed():
		m.contextStatus = StatusFailed
	case status.IsUnstarted():
		m.contextStatus = StatusNotStarted
	case status.IsStarted():
		m.contextStatus = StatusStarted
	case status.IsStopped():
		m.contextStatus = StatusStopped
	default:
		m.contextStatus = StatusUnknown
	}
}

func (m *Manager) buildContextInfo(contextName string) (Detail, error) {
	if m.err != nil {
		return Detail{}, m.err
	}
	contextInfo := Detail{
		Summary: Summary{
			Name:          contextName,
			IsSpot:        m.projectSpec.Contexts[contextName].RequestSpotInstances,
			MaxVCpus:      m.projectSpec.Contexts[contextName].MaxVCpus,
			InstanceTypes: m.projectSpec.Contexts[contextName].InstanceTypes,
		},
		Status:             m.contextStatus,
		BucketLocation:     s3.RenderS3Uri(m.outputBucket, awsresources.RenderBucketContextKey(m.projectSpec.Name, m.userId, contextName)),
		WesUrl:             m.contextStackInfo.Outputs["WesUrl"],
		WesLogGroupName:    m.contextStackInfo.Outputs["AdapterLogGroupName"],
		EngineLogGroupName: m.contextStackInfo.Outputs["EngineLogGroupName"],
		AccessLogGroupName: m.contextStackInfo.Outputs["AccessLogGroupName"],
	}
	return contextInfo, m.err
}
