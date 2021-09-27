package cdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

func ParseOutput(outputPath string) (map[string]string, error) {
	log.Debug().Msgf("ParseOutput(%s)", outputPath)
	jsonFile, err := os.Open(outputPath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	svcOutput := make(map[string]interface{})
	err = json.Unmarshal(byteValue, &svcOutput)
	if err != nil {
		return nil, err
	}

	var outputValues = make(map[string]string)
	err = convertToMap(outputValues, svcOutput)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("ReadJsonFile outputValues: %v", outputValues)
	return outputValues, nil
}

func convertToMap(outputValuesMap map[string]string, svcOutput map[string]interface{}) error {
	for key, val := range svcOutput {
		switch valType := val.(type) {
		case map[string]interface{}:
			err := convertToMap(outputValuesMap, valType)
			if err != nil {
				return err
			}
		case string:
			outputValuesMap[key] = valType
		case float64:
			outputValuesMap[key] = fmt.Sprintf("%f", valType)
		default:
			return fmt.Errorf(key, "is of unexpected type")
		}
	}
	return nil
}
