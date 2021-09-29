package wesmocks

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes"
)

type WesClient interface {
	wes.Interface
}
