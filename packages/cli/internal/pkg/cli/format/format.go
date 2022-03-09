package format

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

var Default Formatter = NewText()

type FormatterType string

const (
	textFormat    FormatterType = "text"
	tableFormat   FormatterType = "table"
	jsonFormat    FormatterType = "json"
	DefaultFormat               = textFormat
)

func NewText() *Text {
	return &Text{os.Stdout}
}

func NewTable(output io.Writer) *Table {
	return &Table{
		*tabwriter.NewWriter(output, 0, 8, 0, tableDelimiter[0], 0),
	}
}

func NewJson() *Json {
	return &Json{os.Stdout}
}

func SetFormatter(format FormatterType) {
	switch format {
	case textFormat:
		Default = NewText()
	case tableFormat:
		Default = NewTable(os.Stdout)
	case jsonFormat:
		Default = NewJson()
	}
}

type Formatter interface {
	Write(interface{})
}

func NewStringFormatter(buffer *bytes.Buffer) Formatter {
	return &Text{buffer}
}

func (f FormatterType) ValidateFormatter() error {
	switch f {
	case textFormat, tableFormat, jsonFormat:
		return nil
	}
	return fmt.Errorf("invalid format type. Valid format types are 'text', 'table', or 'json'")
}
