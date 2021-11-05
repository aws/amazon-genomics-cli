package awsresources

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func GetContextStackStatus(cfn cfn.Interface, projectName string, userId string, contextName string) (types.StackStatus, error) {
	engineStackName := RenderContextStackName(projectName, contextName, userId)
	status, err := cfn.GetStackStatus(engineStackName)
	if err != nil {
		return "", err
	}
	return status, nil
}

func GetCoreStackStatus(cfn cfn.Interface) (types.StackStatus, error) {
	stackName := RenderCoreStackName()
	status, err := cfn.GetStackStatus(stackName)
	if err != nil {
		return "", err
	}
	return status, nil
}
