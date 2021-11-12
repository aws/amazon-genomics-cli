package storagemocks

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
)

type ProjectClient interface {
	storage.ProjectClient
}
type ConfigClient interface {
	config.ConfigClient
}
type StorageClient interface {
	storage.StorageClient
}
