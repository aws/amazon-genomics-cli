package cfn

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

type Client struct {
	Interface
	cfn cfnInterface
}

func New(cfg aws.Config) *Client {
	return &Client{
		cfn: cloudformation.NewFromConfig(cfg),
	}
}
