package iomocks

type OS interface {
	MkdirTemp(dir, pattern string) (string, error)
	RemoveAll(path string) error
}
