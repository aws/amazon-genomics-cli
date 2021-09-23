package format

import "os"

var Default Formatter = &Text{os.Stdout}

type Formatter interface {
	Write(interface{})
}
