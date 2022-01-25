package storagemocks

import "github.com/aws/amazon-genomics-cli/internal/pkg/storage"

type ProjectClient interface {
	storage.ProjectClient
}
type ConfigClient interface {
	storage.ConfigClient
}
type StorageClient interface {
	storage.StorageClient
}
type InputClient interface {
	storage.InputClient
}
