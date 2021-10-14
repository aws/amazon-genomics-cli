package version

import (
	"context"
	"strings"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/environment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

const (
	ChannelVarName = "AGC_UPDATE_CHANNEL"
	DefaultChannel = "s3://healthai-public-assets-us-east-1/amazon-genomics-cli/"

	UpdateCheckCtrlVarName = "AGC_UPDATE_NOTIFICATION"
	DefaultUpdateCheckCtrl = "true"
)

var (
	newS3ClientFromConfig = func(cfg aws.Config) S3Api {
		return s3.NewFromConfig(cfg)
	}

	checkVersion = func(s3Client S3Api, channel string, currentTime time.Time) (Result, error) {
		replenishFunc := newReplenishFromS3Func(s3Client, channel)
		store := &cachedStore{replenishFunc}
		checker := &checker{store, currentTime}
		return checker.Check(Version)
	}

	getCurrentTime = func() time.Time {
		return time.Now()
	}
)

func Check() (Result, error) {
	updateCheckCtrl := environment.LookUpEnvOrDefault(UpdateCheckCtrlVarName, DefaultUpdateCheckCtrl)
	if shouldSkipUpdateCheck(updateCheckCtrl) {
		log.Warn().Msgf("AGC version check is disabled. To re-enable version check unset environment variable '%s'", UpdateCheckCtrlVarName)
		return Result{
			CurrentVersion: Version,
			LatestVersion:  Version,
		}, nil
	}

	channel := environment.LookUpEnvOrDefault(ChannelVarName, DefaultChannel)
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return Result{}, err
	}
	s3Client := newS3ClientFromConfig(cfg)
	currentTime := getCurrentTime()
	return checkVersion(s3Client, channel, currentTime)
}

func shouldSkipUpdateCheck(ctrl string) bool {
	switch strings.ToLower(strings.TrimSpace(ctrl)) {
	case "off", "false", "stop", "0":
		return true
	default:
		return false
	}
}
