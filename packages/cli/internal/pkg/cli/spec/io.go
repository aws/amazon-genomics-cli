package spec

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

var readFile = ioutil.ReadFile
var jsonUnmarshal = json.Unmarshal

//go:embed project_schema.json
var projectSchema string

func ToYaml(specFilePath string, projectSpec Project) error {
	bytes, err := yaml.Marshal(projectSpec)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(specFilePath, bytes, 0644)
}

func FromYaml(specFilePath string) (Project, error) {
	bytes, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return Project{}, err
	}

	if err := ValidateProject(bytes); err != nil {
		return Project{}, err
	}

	var projectSpec Project
	if err := yaml.Unmarshal(bytes, &projectSpec); err != nil {
		return Project{}, err
	}
	return projectSpec, nil
}

func FromJson(manifestFilePath string) (Manifest, error) {
	bytes, err := readFile(manifestFilePath)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err := jsonUnmarshal(bytes, &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func ValidateProject(yamlBytes []byte) error {

	schemaLoader := gojsonschema.NewStringLoader(projectSchema)

	var data interface{}
	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		return err
	}
	structLoader := gojsonschema.NewGoLoader(convertDocumentNode(data))

	result, err := gojsonschema.Validate(schemaLoader, structLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return projectSpecValidationError(result.Errors())
	}

	return nil
}

func projectSpecValidationError(errors []gojsonschema.ResultError) error {
	var errBuffer bytes.Buffer
	errBuffer.WriteString("\n")
	for idx, desc := range errors {
		errBuffer.WriteString(fmt.Sprintf("\t%d. %s\n", idx+1, desc))
	}
	return fmt.Errorf(errBuffer.String())
}

// convertDocumentNode converts yaml derived interfaces into map[string]interface{}
func convertDocumentNode(val interface{}) interface{} {
	if listValue, ok := val.([]interface{}); ok {
		res := []interface{}{}
		for _, v := range listValue {
			res = append(res, convertDocumentNode(v))
		}
		return res
	}
	if mapValue, ok := val.(map[interface{}]interface{}); ok {
		res := map[string]interface{}{}
		for k, v := range mapValue {
			res[k.(string)] = convertDocumentNode(v)
		}
		return res
	}
	return val
}
