package sts

import (
	"github.com/stretchr/testify/mock"
)

type StsMock struct {
	stsInterface
	mock.Mock
}

func NewMockClient() Client {
	return Client{
		sts: new(StsMock),
	}
}
