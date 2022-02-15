package workflow

import (
	"encoding/json"
)

type optionFileKey = string
type optionFileUrl = interface{}
type OptionFile map[optionFileKey]optionFileUrl

func (o OptionFile) ToString() (string, error) {
	jsonBytes, err := json.Marshal(o)
	return string(jsonBytes), err
}
