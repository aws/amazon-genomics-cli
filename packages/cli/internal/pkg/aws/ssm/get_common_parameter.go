package ssm

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

const (
	parameterPrefix       = "/agc/_common"
	outputBucketParameter = "bucket"
	customTagsParameter   = "customTags"
)

func (c *Client) GetCommonParameter(parameterSuffix string) (string, error) {
	parameterName := path.Join(parameterPrefix, parameterSuffix)
	input := &ssm.GetParameterInput{
		Name: aws.String(parameterName),
	}
	output, err := c.ssm.GetParameter(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("unable to obtain bucket output name. The SSM parameter %s may be misconfigured. Error is: %s", parameterName, err)
	}
	if output.Parameter.Value == nil {
		return "", fmt.Errorf("parameter '%s' is not set", parameterName)
	}
	return aws.ToString(output.Parameter.Value), nil
}

func (c *Client) GetOutputBucket() (string, error) {
	return c.GetCommonParameter(outputBucketParameter)
}

func (c *Client) GetCustomTags() string {
	// Custom tags may not exist, so ignore the error

	tags, _ := c.GetCommonParameter(customTagsParameter)
	return tags
}
