package context

type Interface interface {
	Deploy(contextName string, showProgress bool) error
	Info(contextName string) (Detail, error)
	List() (map[string]Summary, error)
	StatusList() ([]Instance, error)
	Destroy(contextName string, showProgress bool) error
}
