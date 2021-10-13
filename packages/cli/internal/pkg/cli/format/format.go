package format

import (
	"bytes"
	"os"
)

var Default Formatter = &Text{os.Stdout}

type Formatter interface {
	Write(interface{})
}

func NewStringFormatter(buffer *bytes.Buffer) Formatter {
	return &Text{buffer}
}
