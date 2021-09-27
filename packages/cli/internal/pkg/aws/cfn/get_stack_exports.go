package cfn

func (c Client) GetStackOutputs(stackName string) (map[string]string, error) {
	info, err := c.GetStackInfo(stackName)
	return info.Outputs, err
}
