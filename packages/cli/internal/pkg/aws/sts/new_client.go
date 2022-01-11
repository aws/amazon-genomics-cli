package sts

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Interface interface {
	GetAccount() (string, error)
}

type stsInterface interface {
	GetCallerIdentity(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type Client struct {
	Interface
	sts stsInterface
}

func NewClient(cfg aws.Config) Interface {
	return Client{
		sts: sts.NewFromConfig(cfg),
	}
}
