package ecr

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type Client struct {
	ecr ecrInterface
}

func New(cfg aws.Config) *Client {
	return &Client{
		ecr: ecr.NewFromConfig(cfg),
	}
}
