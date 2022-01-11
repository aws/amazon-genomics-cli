package ecr

import (
	"context"
	"testing"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func TestClient_ImageListable(t *testing.T) {

	type args struct {
		registry       string
		region         string
		repositoryName string
		imageTag       string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ListableImage",
			args: args{
				registry:       "123456789123",
				region:         "us-east-1",
				repositoryName: "foo",
				imageTag:       "latest",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "NotListableImage",
			args: args{
				registry:       "123456789123",
				region:         "us-east-1",
				repositoryName: "not-a-real-repository",
				imageTag:       "latest",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMockClient()

			call := c.ecr.(*EcrMock).On("ListImages", context.Background(), &ecr.ListImagesInput{
				RepositoryName: &tt.args.repositoryName,
				RegistryId:     &tt.args.registry,
			})

			if tt.args.repositoryName == "foo" {
				call.Return(&ecr.ListImagesOutput{
					ImageIds: []types.ImageIdentifier{{ImageDigest: aws.String("1234eads"), ImageTag: aws.String("latest")}},
				}, nil)
			} else {
				call.Return(&ecr.ListImagesOutput{ImageIds: []types.ImageIdentifier{}}, nil)
			}

			got, err := c.ImageListable(tt.args.registry, tt.args.repositoryName, tt.args.imageTag, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageListable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImageListable() got = %v, want %v", got, tt.want)
			}
		})
	}
}
