package cfn

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func (c Client) GetStackStatus(stackName string) (types.StackStatus, error) {
	info, err := c.GetStackInfo(stackName)
	return info.Status, err
}
