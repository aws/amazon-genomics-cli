package workflow

import (
	"encoding/json"
)

type inputKey = string
type inputUrl = interface{}
type Input map[inputKey]inputUrl

func (i Input) ToString() (string, error) {
	jsonBytes, err := json.Marshal(i)
	return string(jsonBytes), err
}
