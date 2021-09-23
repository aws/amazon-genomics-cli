package cwl

import (
	"github.com/stretchr/testify/mock"
)

type CwlMock struct {
	cwlInterface
	mock.Mock
}

func NewMockClient() Client {
	return Client{
		cwl: new(CwlMock),
	}
}
