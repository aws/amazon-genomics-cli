package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3 s3Interface
}

func New(cfg aws.Config) *Client {
	return &Client{
		s3: s3.NewFromConfig(cfg),
	}
}
