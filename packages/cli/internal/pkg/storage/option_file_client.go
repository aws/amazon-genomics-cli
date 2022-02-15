package storage

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/rs/zerolog/log"
)

type OptionInstance struct {
	S3 s3.Interface
}

func NewOptionClient(S3 s3.Interface) *OptionInstance {
	return &OptionInstance{S3}
}

func (oc *OptionInstance) UpdateOptionReferenceAndUploadToS3(initialProjectDirectory string, tempProjectDirectory string, bucketName string, baseS3Key string) error {
	if !DoesManifestExistInDirectory(tempProjectDirectory) {
		log.Debug().Msgf("No Manifest file was found in folder %s", tempProjectDirectory)
		return nil
	}

	manifest, err := ReadManifestInDirectory(tempProjectDirectory)
	if err != nil {
		return err
	}

	fileLocation := fmt.Sprintf("%s/%s", tempProjectDirectory, manifest.OptionFileUrl)
	optionsReferenceFile, err := ioutilReadFile(fileLocation)
	if err != nil {
		return err
	}

	var optionFile map[string]interface{}
	err = jsonUnmarshall(optionsReferenceFile, &optionFile)
	if err != nil {
		return actionableerror.New(err, fmt.Sprintf("Please validate that the options JSON file %s exists", manifest.OptionFileUrl))
	}

	_, err = oc.UpdateOptionFile(initialProjectDirectory, optionFile, bucketName, baseS3Key, fileLocation)
	if err != nil {
		return err
	}

	return nil
}

func (oc *OptionInstance) UpdateOptionFile(initialProjectDirectory string, optionFile map[string]interface{}, bucketName string, baseS3Key string, fileLocation string) (map[string]interface{}, error) {
	var updatedOptionReferenceFile = make(map[string]interface{})
	for key, value := range optionFile {
		var optionReference string
		switch typedValue := value.(type) {
		case string:
			updatedReference, err := oc.uploadReferenceToS3(optionReference, initialProjectDirectory, bucketName, baseS3Key)
			if err != nil {
				return nil, err
			}

			updatedOptionReferenceFile[key] = updatedReference
		case []interface{}:
			var updatedRef []interface{}

			// We only support one level deep
			for _, val := range typedValue {
				stringValue, ok := val.(string)
				if !ok {
					updatedRef = append(updatedRef, val)
					log.Debug().Msgf("The value %#v is not a string and will not be checked if it's an options file", val)
				} else {
					optionReference = stringValue
				}
			}
			updatedReferences, err := oc.uploadReferenceToS3(optionReference, initialProjectDirectory, bucketName, baseS3Key)
			if err != nil {
				return nil, err
			}

			for _, val := range updatedReferences {
				updatedRef = append(updatedRef, val)
			}

			updatedOptionReferenceFile[key] = updatedRef
		default:
			updatedOptionReferenceFile[key] = value
		}
	}
	marshalledData, err := jsonMarshall(updatedOptionReferenceFile)
	if err != nil {
		return nil, err
	}
	err = ioutilWriteFile(fileLocation, marshalledData, 0644)
	if err != nil {
		return nil, err
	}

	return updatedOptionReferenceFile, nil
}

func (oc *OptionInstance) uploadReferenceToS3(optionFileLocation string, baseDirectory string, bucketName string, baseS3Key string) (string, error) {
	var updatedReference string
	trimmedInput := strings.TrimSpace(optionFileLocation)
	optionsWithDirectory := fmt.Sprintf("%s/%s", baseDirectory, trimmedInput)
	if _, err := os.Stat(optionsWithDirectory); err == nil {
		var formattedOptionFileName string
		if strings.HasPrefix(trimmedInput, "./") {
			formattedOptionFileName = trimmedInput[2:]
		} else {
			formattedOptionFileName = trimmedInput
		}
		err = oc.S3.UploadFile(bucketName, fmt.Sprintf("%s/%s", baseS3Key, formattedOptionFileName), optionsWithDirectory)
		if err != nil {
			return "", err
		}
		updatedReference = fmt.Sprintf("s3://%s/%s/%s", bucketName, baseS3Key, formattedOptionFileName)
	} else {
		updatedReference = optionFileLocation
		log.Debug().Msgf("The following option file value is not a file %s", err)
	}

	return updatedReference, nil
}
