package ecr

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/stretchr/testify/mock"
)

type EcrMock struct {
	ecr.ListImagesAPIClient
	mock.Mock
}

func NewMockClient() *Client {
	return &Client{
		ecr: new(EcrMock),
	}
}

func (m *EcrMock) ListImages(ctx context.Context, input *ecr.ListImagesInput, opts ...func(*ecr.Options)) (*ecr.ListImagesOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*ecr.ListImagesOutput), err
	}
	return nil, err
}
