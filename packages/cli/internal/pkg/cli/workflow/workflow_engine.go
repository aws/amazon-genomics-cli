package workflow

import "github.com/rsc/wes_client"

type EngineLog struct {
	WorkflowRunId  string
	StdOut         string
	StdErr         string
	WorkflowStatus wes_client.State
}

func (m *Manager) GetEngineLogByRunId(runId string) (EngineLog, error) {
	m.readProjectSpec()
	m.readConfig()
	m.setContextForRun(runId)
	m.setContext(m.runContextName)
	m.setEngineForWorkflowType(m.runContextName)
	m.setContextStackInfo(m.runContextName)
	m.setWesUrl()
	m.setWesClient()
	m.getRunLog(runId)

	return m.buildEngineLog()
}

func (m *Manager) buildEngineLog() (EngineLog, error) {
	if m.err != nil {
		return EngineLog{}, m.err
	}

	return EngineLog{
		WorkflowRunId:  m.taskProps.runLog.RunId,
		StdOut:         m.taskProps.runLog.RunLog.Stdout,
		StdErr:         m.taskProps.runLog.RunLog.Stderr,
		WorkflowStatus: m.taskProps.runLog.State,
	}, nil
}
