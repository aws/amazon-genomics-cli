package types

type Engine struct {
	Name string
	Spec EngineSpec
}

type EngineSpec struct {
	Raw string //TODO: update with actual structure when Project specification if finalized
}

type EngineInstance struct {
	Id       string
	Name     string
	Status   string
	Error    string
	Start    string
	Duration string
	Log      EngineLog
}

type EngineLog struct {
	Raw string
}
