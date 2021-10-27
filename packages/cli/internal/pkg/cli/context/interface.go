package context

type Interface interface {
	Deploy(contexts []string) []ProgressResult
	Info(contextName string) (Detail, error)
	List() (map[string]Summary, error)
	StatusList() ([]Instance, error)
	Destroy(contexts []string) []ProgressResult
}
