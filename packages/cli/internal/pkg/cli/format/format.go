package format

import (
	"os"
	"text/tabwriter"
)
var Default Formatter

type FormatterType string
const (
	textFormat FormatterType = "text"
	tabularFormat FormatterType = "tabular"
	DefaultFormat = textFormat
)

var Formats = map[FormatterType]func() Formatter{
	textFormat:    func() Formatter { return &Text{os.Stdout} },
	tabularFormat: func() Formatter {
		return &TextTabular{
		*tabwriter.NewWriter(os.Stdout, 0, 8, 0, tabularDelimiter[0], 0),
	}},
}

func SetFormatter(format FormatterType) {
	if writer, ok := Formats[format]; ok {
		Default = writer()
	} else {
		Default = Formats[DefaultFormat]()
	}
}

type Formatter interface {
	Write(interface{})
}
