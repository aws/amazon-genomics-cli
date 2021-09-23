package clierror

import (
	"bytes"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type ActionableError struct {
	Cause           error
	SuggestedAction string
}

func (e ActionableError) Error() string {
	return fmt.Sprintf("an error occurred caused by: %s\nsuggestion: %s", e.Cause, e.SuggestedAction)
}

func ProjectSpecValidationError(errors []gojsonschema.ResultError) ActionableError {
	var errBuffer bytes.Buffer
	errBuffer.WriteString("\n")
	for idx, desc := range errors {
		errBuffer.WriteString(fmt.Sprintf("\t%d. %s\n", idx+1, desc))
	}
	return ActionableError{
		Cause:           fmt.Errorf(errBuffer.String()),
		SuggestedAction: "please fix the validation errors in the AGC project file",
	}
}
