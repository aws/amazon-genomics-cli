package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ioutilReadFile  = ioutil.ReadFile
	ioutilWriteFile = ioutil.WriteFile
	jsonUnmarshall  = json.Unmarshal
	jsonMarshall    = json.Marshal
)

func (ic *InputInstance) UpdateInputReferencesAndUploadToS3(initialProjectDirectory string, tempProjectDirectory string, bucketName string, baseS3Key string) error {
	if !DoesManifestExistInDirectory(tempProjectDirectory) {
		log.Debug().Msgf("No Manifest file was found in folder %s", tempProjectDirectory)
		return nil
	}

	manifest, err := ReadManifestInDirectory(tempProjectDirectory)
	if err != nil {
		return err
	}

	for _, inputLocation := range manifest.InputFileUrls {
		fileLocation := fmt.Sprintf("%s/%s", tempProjectDirectory, inputLocation)
		inputReferenceFile, err := ioutilReadFile(fileLocation)
		if err != nil {
			return err
		}

		var inputFile map[string]interface{}
		err = jsonUnmarshall(inputReferenceFile, &inputFile)
		if err != nil {
			return actionableerror.New(err, fmt.Sprintf("Please validate that the input JSON file %s exists", inputLocation))
		}

		err = ic.updateInputsInFile(initialProjectDirectory, inputFile, bucketName, baseS3Key, fileLocation)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ic *InputInstance) updateInputsInFile(initialProjectDirectory string, inputFile map[string]interface{}, bucketName string, baseS3Key string, fileLocation string) error {
	var updatedInputReferenceFile = make(map[string]interface{})
	for key, value := range inputFile {
		var inputReferences []string
		switch typedValue := value.(type) {
		case string:
			inputReferences = strings.Split(typedValue, ",")
			updatedReferences, err := ic.uploadReferencesToS3(inputReferences, initialProjectDirectory, bucketName, baseS3Key)
			if err != nil {
				return err
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
				return err
			}

			for _, val := range updatedReferences {
				updatedRef = append(updatedRef, val)
			}

			updatedInputReferenceFile[key] = updatedRef
		default:
			updatedInputReferenceFile[key] = value
		}
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
		if _, err := os.Stat(inputWithDirectory); err == nil {
			var formattedInputName string
			if strings.HasPrefix(trimmedInput, "./") {
				formattedInputName = trimmedInput[2:]
			} else {
				formattedInputName = trimmedInput
			}
			err = ic.S3.UploadFile(bucketName, fmt.Sprintf("%s/%s", baseS3Key, formattedInputName), inputWithDirectory)
			if err != nil {
				return nil, err
			}
			updatedReferences[index] = fmt.Sprintf("s3://%s/%s/%s", bucketName, baseS3Key, formattedInputName)
		} else {
			updatedReferences[index] = input
			log.Debug().Msgf("The following input value is not a file %s", err)
		}
	}

	return updatedReferences, nil
}
