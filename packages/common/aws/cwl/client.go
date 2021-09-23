package cwl

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type Client struct {
	Interface
	cwl cwlInterface
}

func New(cfg aws.Config) *Client {
	return &Client{
		cwl: cloudwatchlogs.NewFromConfig(cfg),
	}
}
