package workflow

import (
	"encoding/json"
)

type Input map[string]interface{}

func (i Input) String() string {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
