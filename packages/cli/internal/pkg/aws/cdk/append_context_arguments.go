package cdk

func appendContextArguments(arguments []string, contextArguments []string) []string {
	for _, envVar := range contextArguments {
		arguments = append(arguments, "-c", envVar)
	}
	return arguments
}
