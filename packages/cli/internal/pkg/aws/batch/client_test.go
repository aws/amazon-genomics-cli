package batch

import "github.com/stretchr/testify/mock"
type BatchMock struct {
	batchInterface
	mock.Mock
}
func NewMockClient() Client {
	return Client{
		batch: new(BatchMock),
	}
}
