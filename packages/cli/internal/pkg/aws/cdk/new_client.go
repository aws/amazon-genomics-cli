package cdk

type Interface interface {
	ClearContext(appDir string) error
	DeployApp(appDir string, context []string) (ProgressStream, error)
	DestroyApp(appDir string, context []string) (ProgressStream, error)
}

type Client struct {
	Interface
	profile string
}

func NewClient(profile string) Interface {
	return Client{
		profile: profile,
	}
}
