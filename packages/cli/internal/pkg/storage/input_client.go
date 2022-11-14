package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/rs/zerolog/log"
)

type InputInstance struct {
	S3 s3.Interface
}

func NewInputClient(S3 s3.Interface) *InputInstance {
	return &InputInstance{S3}
}

var (
	ioutilReadFile  = os.ReadFile
	ioutilWriteFile = os.WriteFile
	jsonUnmarshall  = json.Unmarshal
	jsonMarshall    = json.Marshal
	stat            = os.Stat
)

func (ic *InputInstance) UpdateInputReferencesAndUploadToS3(initialProjectDirectory string, tempProjectDirectory string, bucketName string, baseS3Key string) error {
	if !DoesManifestExistInDirectory(tempProjectDirectory) {
		log.Debug().Msgf("no Manifest file was found in folder %s, input references will not be updated to s3 locations", tempProjectDirectory)
		return nil
	}

	log.Debug().Msgf("reading manifest in '%s", tempProjectDirectory)
	manifest, err := ReadManifestInDirectory(tempProjectDirectory)
	if err != nil {
		return err
	}

	log.Debug().Msgf("manifest declares '%d' input files", len(manifest.InputFileUrls))
	for _, inputLocation := range manifest.InputFileUrls {
		fileLocation := fmt.Sprintf("%s/%s", tempProjectDirectory, inputLocation)
		log.Debug().Msgf("reading content of input file at '%s'", fileLocation)
		inputReferenceFile, err := ioutilReadFile(fileLocation)
		if err != nil {
			return err
		}
		log.Debug().Msgf("content of '%s' is \n%s", fileLocation, string(inputReferenceFile))

		var inputFile map[string]interface{}
		err = jsonUnmarshall(inputReferenceFile, &inputFile)
		if err != nil {
			return actionableerror.New(err, fmt.Sprintf("Please validate that the input JSON file %s exists and that the content is valid JSON", inputLocation))
		}

		err = ic.updateInputsInFile(initialProjectDirectory, inputFile, bucketName, baseS3Key, fileLocation)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ic *InputInstance) UpdateInputs(initialProjectDirectory string, inputFile map[string]interface{}, bucketName string, baseS3Key string) (map[string]interface{}, error) {
	var updatedInputReferenceFile = make(map[string]interface{})
	for key, value := range inputFile {
		log.Debug().Msgf("inspecting key value pair, '%s: %s'", key, value)
		var inputReferences []string
		switch typedValue := value.(type) {
		case string:
			inputReferences = strings.Split(typedValue, ",")
			updatedReferences, err := ic.uploadReferencesToS3(inputReferences, initialProjectDirectory, bucketName, baseS3Key)
			if err != nil {
				return nil, err
			}

			updatedInputReferenceFile[key] = strings.Join(updatedReferences, ",")
		case []interface{}:
			var updatedRef []interface{}

			// We only support one level deep
			for _, val := range typedValue {
				stringValue, ok := val.(string)
				if !ok {
					updatedRef = append(updatedRef, val)
					log.Debug().Msgf("The value %#v is not a string and will not be checked if it's an input file", val)
				} else {
					inputReferences = append(inputReferences, stringValue)
				}
			}
			updatedReferences, err := ic.uploadReferencesToS3(inputReferences, initialProjectDirectory, bucketName, baseS3Key)
			if err != nil {
				return nil, err
			}

			for _, val := range updatedReferences {
				updatedRef = append(updatedRef, val)
			}

			updatedInputReferenceFile[key] = updatedRef
		default:
			updatedInputReferenceFile[key] = value
		}
		log.Debug().Msgf("key value pair updated to '%s: %s'", key, updatedInputReferenceFile[key])
	}

	return updatedInputReferenceFile, nil
}

func (ic *InputInstance) updateInputsInFile(initialProjectDirectory string, inputFile map[string]interface{}, bucketName string, baseS3Key string, fileLocation string) error {
	updatedInputReferenceFile, err := ic.UpdateInputs(initialProjectDirectory, inputFile, bucketName, baseS3Key)
	if err != nil {
		return err
	}
	marshalledData, err := jsonMarshall(updatedInputReferenceFile)
	if err != nil {
		return err
	}
	err = ioutilWriteFile(fileLocation, marshalledData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ic *InputInstance) uploadReferencesToS3(inputLocations []string, baseDirectory string, bucketName string, baseS3Key string) ([]string, error) {
	var updatedReferences = make([]string, len(inputLocations))
	for index, input := range inputLocations {
		trimmedInput := strings.TrimSpace(input)
		inputWithDirectory := fmt.Sprintf("%s/%s", baseDirectory, trimmedInput)
		if fInfo, err := stat(inputWithDirectory); err == nil && !fInfo.IsDir() {
			log.Debug().Msgf("input value '%s' can be resolved to a file at '%s'", trimmedInput, inputWithDirectory)
			var formattedInputName string
			if strings.HasPrefix(trimmedInput, "./") {
				formattedInputName = trimmedInput[2:]
			} else {
				formattedInputName = trimmedInput
			}
			s3Location := fmt.Sprintf("%s/%s", baseS3Key, formattedInputName)
			log.Debug().Msgf("loading '%s' to '%s'", formattedInputName, s3Location)
			err = ic.S3.UploadFile(bucketName, s3Location, inputWithDirectory)
			if err != nil {
				return nil, err
			}

			s3Reference := fmt.Sprintf("s3://%s/%s/%s", bucketName, baseS3Key, formattedInputName)
			updatedReferences[index] = s3Reference
			log.Debug().Msgf("updated reference '%d' to '%s'", index, s3Reference)
		} else {
			updatedReferences[index] = input
			log.Debug().Msgf("The following input value is not a file %s", err)
		}
	}

	return updatedReferences, nil
}
