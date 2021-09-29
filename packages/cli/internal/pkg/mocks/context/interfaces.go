package contextmocks

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
)

type ContextManager interface {
	context.Interface
}
