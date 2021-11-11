package version

import (
	"context"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/blang/semver/v4"
)

const (
	deprecationMessageObjectName = "deprecated"
	highlightMessageObjectName   = "highlight"
)

type semanticVersion struct {
	Version semver.Version
	Info
}

func newReplenishFromS3Func(s3Client S3Api, channel string) func(string) ([]Info, error) {
	return func(versionString string) ([]Info, error) {
		currentSemVersion, err := semver.Parse(versionString)
		if err != nil {
			return nil, err
		}
		bucket, prefix, err := parseS3Url(channel)
		if err != nil {
			return nil, err
		}
		prefix = stripLeadingSlashIfAny(prefix)
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		}
		index := make(map[string]*semanticVersion)
		var semanticVersions []*semanticVersion
		paginator := s3.NewListObjectsV2Paginator(s3Client, input)
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.Background())
			if err != nil {
				return nil, err
			}
			s3Objects := output.Contents

			for _, s3Object := range s3Objects {
				// Skipping all unexpected objects in the bucket
				// We are looking only for:
				// - <prefix>/<semantic version>/deprecated
				// - <prefix>/<semantic version>/highlight

				key := aws.ToString(s3Object.Key)
				if !strings.HasPrefix(key, prefix) {
					continue
				}
				versionAndObjectName := strings.Split(key[len(prefix):], "/")
				if len(versionAndObjectName) != 2 {
					continue
				}
				versionString := versionAndObjectName[0]
				versionMeta, ok := index[versionString]
				if !ok {
					semVersion, err := semver.Parse(versionString)
					if err != nil {
						continue
					}
					versionMeta = &semanticVersion{Version: semVersion}
					versionMeta.Name = versionString
					semanticVersions = append(semanticVersions, versionMeta)
					index[versionString] = versionMeta
				}
				objectName := versionAndObjectName[1]
				if currentSemVersion.EQ(versionMeta.Version) && objectName == deprecationMessageObjectName {
					versionMeta.Deprecated = true
					versionMeta.DeprecationMessage, _ = readS3Object(s3Client, bucket, key)
					continue
				}
				if currentSemVersion.LT(versionMeta.Version) && objectName == highlightMessageObjectName {
					versionMeta.Highlight, _ = readS3Object(s3Client, bucket, key)
					continue
				}
			}
		}
		sort.Slice(semanticVersions, func(i, j int) bool {
			return semanticVersions[i].Version.LT(semanticVersions[j].Version)
		})
		currentVersionIndex := sort.Search(len(semanticVersions), func(i int) bool {
			return semanticVersions[i].Version.GTE(currentSemVersion)
		})
		currentAndNewer := semanticVersions[currentVersionIndex:]

		result := make([]Info, len(currentAndNewer))
		for i, ver := range currentAndNewer {
			result[i] = ver.Info
		}
		return result, nil
	}
}

func stripLeadingSlashIfAny(path string) string {
	if strings.HasPrefix(path, "/") {
		return path[1:]
	}
	return path
}

func parseS3Url(channelUrl string) (string, string, error) {
	parsedUrl, err := url.Parse(channelUrl)
	if err != nil {
		return "", "", err
	}
	return parsedUrl.Host, parsedUrl.Path, nil
}

func readS3Object(s3Client S3Api, bucket string, key string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	output, err := s3Client.GetObject(context.Background(), input)
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(output.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
