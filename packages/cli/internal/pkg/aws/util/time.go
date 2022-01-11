package util

import (
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func TimeToAws(someTime *time.Time) *int64 {
	if someTime == nil {
		return nil
	}
	return aws.Int64(nanoToMilli(someTime.UnixNano()))
}

func TimeFromAws(someTime *int64) time.Time {
	return time.Unix(0, milliToNano(aws.ToInt64(someTime)))
}

func nanoToMilli(nano int64) int64 {
	return nano / 1000000
}

func milliToNano(milli int64) int64 {
	return milli * 1000000
}
