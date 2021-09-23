package cfn

import (
	"github.com/stretchr/testify/mock"
)

type CfnMock struct {
	cfnInterface
	mock.Mock
}

func NewMockClient() Client {
	return Client{
		cfn: new(CfnMock),
	}
}
