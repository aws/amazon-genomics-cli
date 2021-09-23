package s3

import (
	"fmt"
	"net/url"
)

const s3Scheme = "s3"

func IsS3Uri(value string) bool {
	urlParts, err := url.Parse(value)
	if err != nil {
		return false
	}
	return s3Scheme == urlParts.Scheme
}

func UriToArn(uri string) (string, error) {
	urlParts, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	if s3Scheme != urlParts.Scheme {
		return "", fmt.Errorf("expected an S3 URI but got '%s'", urlParts.Scheme)
	}
	return fmt.Sprintf("arn:aws:s3:::%s%s", urlParts.Host, urlParts.Path), nil
}

func RenderS3Uri(bucketName, objectKey string) string {
	uri := url.URL{
		Scheme: s3Scheme,
		Host:   bucketName,
		Path:   objectKey,
	}
	return uri.String()
}
