package storagemocks

import "github.com/aws/amazon-genomics-cli/cli/internal/pkg/storage"

type ProjectClient interface {
	storage.ProjectClient
}
type ConfigClient interface {
	storage.ConfigClient
}
type StorageClient interface {
	storage.StorageClient
}
