package version

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ReplenishFuncTestSuite struct {
	suite.Suite

	channel        string
	currentVersion string
	prefix         string
	bucket         string

	ctrl          *gomock.Controller
	mockS3Client  *MockS3Api
	replenishFunc func(versionString string) ([]Info, error)
}

func (s *ReplenishFuncTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.channel = DefaultChannel
	s.currentVersion = "1.0.1"
	s.bucket = "healthai-public-assets-us-east-1"
	s.prefix = "amazon-genomics-cli/"

	s.mockS3Client = NewMockS3Api(s.ctrl)
	s.replenishFunc = newReplenishFromS3Func(s.mockS3Client, s.channel)
}

func (s *ReplenishFuncTestSuite) AfterTest(_, _ string) {
	s.ctrl.Finish()
}

func (s *ReplenishFuncTestSuite) TestReplenishNominal() {
	s.mockS3Client.EXPECT().ListObjectsV2(context.Background(), mock.MatchedBy(newListObjectV2InputMatcher(s.bucket, s.prefix))).Return(
		&s3.ListObjectsV2Output{
			KeyCount:    2,
			IsTruncated: false,
			Contents: []types.Object{
				{Key: aws.String("amazon-genomics-cli/")},
				{Key: aws.String("amazon-genomics-cli/1.0.1/")},
			},
		}, nil)

	expected := []Info{
		{Name: s.currentVersion},
	}
	actual, err := s.replenishFunc(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *ReplenishFuncTestSuite) TestReplenishIgnoreUnknownObjects() {
	s.mockS3Client.EXPECT().ListObjectsV2(context.Background(), mock.MatchedBy(newListObjectV2InputMatcher(s.bucket, s.prefix))).Return(
		&s3.ListObjectsV2Output{
			KeyCount:    4,
			IsTruncated: false,
			Contents: []types.Object{
				{Key: aws.String("amazon-genomics-cli/")},
				{Key: aws.String("amazon-genomics-cli/1.0.1/version")},
				{Key: aws.String("amazon-genomics-cli/1.0.1/notes")},
				{Key: aws.String("amazon-genomics-cli/1.0.1/authors")},
			},
		}, nil)

	expected := []Info{
		{Name: s.currentVersion},
	}
	actual, err := s.replenishFunc(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *ReplenishFuncTestSuite) TestReplenishIgnoreInvalidVersions() {
	s.mockS3Client.EXPECT().ListObjectsV2(context.Background(), mock.MatchedBy(newListObjectV2InputMatcher(s.bucket, s.prefix))).Return(
		&s3.ListObjectsV2Output{
			KeyCount:    5,
			IsTruncated: false,
			Contents: []types.Object{
				{Key: aws.String("amazon-genomics-cli/")},
				{Key: aws.String("amazon-genomics-cli/demo/")},
				{Key: aws.String("amazon-genomics-cli/1.0.1/authors")},
				{Key: aws.String("amazon-genomics-cli/beta/")},
				{Key: aws.String("amazon-genomics-cli/poc/")},
			},
		}, nil)

	expected := []Info{
		{Name: s.currentVersion},
	}
	actual, err := s.replenishFunc(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *ReplenishFuncTestSuite) TestReplenishDeprecation() {
	deprecatedKey := "amazon-genomics-cli/1.0.1/deprecated"
	deprecationMessage := "Test Deprecation Message"
	s.mockS3Client.EXPECT().ListObjectsV2(context.Background(), mock.MatchedBy(newListObjectV2InputMatcher(s.bucket, s.prefix))).Return(
		&s3.ListObjectsV2Output{
			KeyCount:    2,
			IsTruncated: false,
			Contents: []types.Object{
				{Key: aws.String("amazon-genomics-cli/")},
				{Key: aws.String(deprecatedKey)},
			},
		}, nil)
	s.mockS3Client.EXPECT().GetObject(context.Background(), mock.MatchedBy(newGetObjectInputMatcher(s.bucket, deprecatedKey))).Return(
		&s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(deprecationMessage))}, nil)

	expected := []Info{
		{
			Name:               s.currentVersion,
			Deprecated:         true,
			DeprecationMessage: deprecationMessage,
		},
	}
	actual, err := s.replenishFunc(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *ReplenishFuncTestSuite) TestReplenishHighlight() {
	highlightKey := "amazon-genomics-cli/1.0.2/highlight"
	highlightMessage := "Test Highlight Message"
	s.mockS3Client.EXPECT().ListObjectsV2(context.Background(), mock.MatchedBy(newListObjectV2InputMatcher(s.bucket, s.prefix))).Return(
		&s3.ListObjectsV2Output{
			KeyCount:    3,
			IsTruncated: false,
			Contents: []types.Object{
				{Key: aws.String("amazon-genomics-cli/")},
				{Key: aws.String("amazon-genomics-cli/1.0.1/")},
				{Key: aws.String(highlightKey)},
			},
		}, nil)
	s.mockS3Client.EXPECT().GetObject(context.Background(), mock.MatchedBy(newGetObjectInputMatcher(s.bucket, highlightKey))).Return(
		&s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(highlightMessage))}, nil)

	expected := []Info{
		{
			Name: s.currentVersion,
		},
		{
			Name:      "1.0.2",
			Highlight: highlightMessage,
		},
	}
	actual, err := s.replenishFunc(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *ReplenishFuncTestSuite) TestReplenishOlderVersions() {
	s.mockS3Client.EXPECT().ListObjectsV2(context.Background(), mock.MatchedBy(newListObjectV2InputMatcher(s.bucket, s.prefix))).Return(
		&s3.ListObjectsV2Output{
			KeyCount:    4,
			IsTruncated: false,
			Contents: []types.Object{
				{Key: aws.String("amazon-genomics-cli/")},
				{Key: aws.String("amazon-genomics-cli/0.99.999/")},
				{Key: aws.String("amazon-genomics-cli/1.0.0/")},
				{Key: aws.String("amazon-genomics-cli/1.0.1/")},
			},
		}, nil)

	expected := []Info{
		{
			Name: s.currentVersion,
		},
	}
	actual, err := s.replenishFunc(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func TestReplenishFuncTestSuite(t *testing.T) {
	suite.Run(t, new(ReplenishFuncTestSuite))
}

func newListObjectV2InputMatcher(bucket, prefix string) func(*s3.ListObjectsV2Input) bool {
	return func(input *s3.ListObjectsV2Input) bool {
		return aws.ToString(input.Bucket) == bucket && aws.ToString(input.Prefix) == prefix
	}
}

func newGetObjectInputMatcher(bucket, key string) func(input *s3.GetObjectInput) bool {
	return func(input *s3.GetObjectInput) bool {
		return aws.ToString(input.Bucket) == bucket && aws.ToString(input.Key) == key
	}
}
