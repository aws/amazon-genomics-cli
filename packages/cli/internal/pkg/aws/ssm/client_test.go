package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/stretchr/testify/mock"
)

type ssmMockClient struct {
	mock.Mock
}

func (s *ssmMockClient) GetParameter(ctx context.Context, input *ssm.GetParameterInput, _ ...func(options *ssm.Options)) (*ssm.GetParameterOutput, error) {
	args := s.Called(ctx, input)
	return args.Get(0).(*ssm.GetParameterOutput), args.Error(1)
}
