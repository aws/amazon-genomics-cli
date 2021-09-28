package ddb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	TableName = "Agc"
	Lsi1Name  = "lsi1"
	Lsi2Name  = "lsi2"
	Lsi3Name  = "lsi3"

	pkAttrName     = "PK"
	skAttrName     = "SK"
	gsi1PkAttrName = "GSI1_PK" //nolint:deadcode,varcheck
	gsi1SkAttrName = "GSI1_SK" //nolint:deadcode,varcheck
	lsi1SkAttrName = "LSI1_SK"
	lsi2SkAttrName = "LSI2_SK"
	lsi3SkAttrName = "LSI3_SK"
)

type Client struct {
	svc ApiInterface
}

func New(cfg aws.Config) *Client {
	return &Client{svc: dynamodb.NewFromConfig(cfg)}
}
