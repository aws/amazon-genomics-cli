package spec

type Manifest struct {
	MainWorkflowUrl string   `json:"mainWorkflowURL,omitempty"`
	InputFileUrls   []string `json:"inputFileURLs,omitempty"`
	OptionFileUrl   string   `json:"optionFileURL,omitempty"`
	EngineOptions   string   `json:"engineOptions,omitempty"`
}
