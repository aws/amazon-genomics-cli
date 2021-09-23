package context

import "fmt"

func environmentMapToList(environmentMap map[string]string) []string {
	var environmentList []string
	for key, value := range environmentMap {
		environmentList = append(environmentList, fmt.Sprintf("%s=%s", key, value))
	}
	return environmentList
}
