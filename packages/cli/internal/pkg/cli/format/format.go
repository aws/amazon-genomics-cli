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

func SetFormatter(format FormatterType) {
	switch format {
	case textFormat:
		Default = NewText()
	case tableFormat:
		Default = NewTable(os.Stdout)
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
	case textFormat, tableFormat:
		return nil
	}
	return fmt.Errorf("invalid format type. Valid format types are 'text' and 'table'")
}
