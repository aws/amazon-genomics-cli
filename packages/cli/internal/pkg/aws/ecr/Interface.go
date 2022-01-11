package ecr

import "github.com/aws/aws-sdk-go-v2/service/ecr"
type Interface interface {
	ImageListable(string, string, string, string) (bool, error)
	VerifyImageExists(reference ImageReference) error
}

type ecrInterface interface {
	ecr.ListImagesAPIClient
}
