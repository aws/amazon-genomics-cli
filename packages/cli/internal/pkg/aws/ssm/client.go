package ssm

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Client struct {
	ssm ssmInterface
}

func New(cfg aws.Config) *Client {
	return &Client{
		ssm: ssm.NewFromConfig(cfg),
	}
}
