package s3

import (
	"github.com/stretchr/testify/mock"
)

type S3Mock struct {
	s3Interface
	mock.Mock
}

func NewMockClient() Client {
	return Client{
		s3: new(S3Mock),
	}
}
