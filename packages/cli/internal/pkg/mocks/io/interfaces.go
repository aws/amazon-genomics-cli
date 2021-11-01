package iomocks

import (
	"io/fs"

	"github.com/rs/zerolog"
)

type OS interface {
	Remove(name string) error
	Chdir(dir string) error
	MkdirTemp(dir, pattern string) (string, error)
	RemoveAll(path string) error
	UserHomeDir() (string, error)
}

type Zip interface {
	CompressToTmp(srcPath string) (string, error)
}

type Tmp interface {
	Write(namePattern, content string) (string, error)
}

type FileReader interface {
	ReadFile(string) ([]byte, error)
}

type FileWriter interface {
	WriteFile(filename string, data []byte, perm fs.FileMode) error
}

type Format interface {
	LogsPrintLn(args ...interface{})
}

type Log interface {
	Info() *zerolog.Event
}
