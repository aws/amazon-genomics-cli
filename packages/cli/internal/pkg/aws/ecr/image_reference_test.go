package ecr

import (
	"context"
	"testing"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

var unpopulatedRef = ImageReference{
	RegistryId:     "",
	Region:         "",
	RepositoryName: "",
	ImageTag:       "",
}

var imageThatExists = ImageReference{
	RegistryId:     "123456788",
	Region:         "region",
	RepositoryName: "repository",
	ImageTag:       "latest",
}

var imageThatDoesntExist = ImageReference{
	RegistryId:     "123456788",
	Region:         "region",
	RepositoryName: "not-a-repository",
	ImageTag:       "latest",
}

func TestVerifyImageExists(t *testing.T) {
	type args struct {
		reference ImageReference
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Unpopulated Reference should produce error", args: args{reference: unpopulatedRef}, wantErr: true},
		{name: "Existing Reference should not error", args: args{reference: imageThatExists}, wantErr: false},
		{name: "Non-existent Reference should error", args: args{reference: imageThatDoesntExist}, wantErr: true},
	}
	for _, tt := range tests {

		c := NewMockClient()

		call := c.ecr.(*EcrMock).On("ListImages", context.Background(), &ecr.ListImagesInput{
			RepositoryName: &tt.args.reference.RepositoryName,
			RegistryId:     &tt.args.reference.RegistryId,
		})

		if tt.args.reference == imageThatExists {
			call.Return(&ecr.ListImagesOutput{
				ImageIds: []types.ImageIdentifier{{ImageDigest: aws.String("1234eads"), ImageTag: aws.String("latest")}},
			}, nil)
		} else {
			call.Return(&ecr.ListImagesOutput{ImageIds: []types.ImageIdentifier{}}, nil)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := c.VerifyImageExists(tt.args.reference); (err != nil) != tt.wantErr {
				t.Errorf("VerifyImageExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
