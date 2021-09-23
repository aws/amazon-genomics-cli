package workflow

import (
	"encoding/json"
)

type inputKey = string
type inputUrl = string
type Input map[inputKey]inputUrl

func (i Input) MapInputUrls(mapper func(key inputKey, fileUrl inputUrl) inputUrl) {
	for key, value := range i {
		i[key] = mapper(key, value)
	}
}

func (i Input) ToString() (string, error) {
	jsonBytes, err := json.Marshal(i)
	return string(jsonBytes), err
}
