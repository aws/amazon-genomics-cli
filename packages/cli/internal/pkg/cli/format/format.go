package format

import (
	"os"
	"text/tabwriter"
)

var Formats = map[string]Formatter{
	"text":    textWriter,
	"tabular": tabularWriter,
}
var Format string

var textWriter Formatter = &Text{os.Stdout}
var tabularWriter = &TextTabular{
	*tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0),
}
var Default Formatter

func SetDefaultWriter() {
	if writer, ok := Formats[Format]; ok {
		Default = writer
	} else {
		Default = textWriter
	}
}

type Formatter interface {
	Write(interface{})
}
