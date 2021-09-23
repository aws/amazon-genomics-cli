package cfn

func (c Client) GetStackTags(stackName string) (map[string]string, error) {
	info, err := c.GetStackInfo(stackName)
	return info.Tags, err
}
