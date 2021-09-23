package contextmocks

import (
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/context"
)

type ContextManager interface {
	context.Interface
}
