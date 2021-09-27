package batch

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
)

type Client struct {
	Interface
	batch batchInterface
}

func New(cfg aws.Config) *Client {
	return &Client{
		batch: batch.NewFromConfig(cfg),
	}
}
