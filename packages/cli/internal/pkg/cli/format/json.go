package format

import (
	"encoding/json"
	"fmt"
	"io"
)

type Json struct {
	writer io.Writer
}

func (f *Json) write(bytes []byte) {
	_, _ = fmt.Fprint(f.writer, string(bytes))
}

func (f *Json) newLine() {
	_, _ = fmt.Fprintln(f.writer)
}

func (f *Json) Write(o interface{}) {
	jsons, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		SetFormatter(DefaultFormat)
		Default.Write(o)
	} else {
		f.write(jsons)
		f.newLine()
	}
}
