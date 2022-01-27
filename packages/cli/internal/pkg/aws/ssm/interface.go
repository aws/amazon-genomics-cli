package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Interface interface {
	GetOutputBucket() (string, error)
	GetCommonParameter(parameterSuffix string) (string, error)
	GetCustomTags() string
}

type ssmInterface interface {
	GetParameter(context.Context, *ssm.GetParameterInput, ...func(options *ssm.Options)) (*ssm.GetParameterOutput, error)
}
