package workflow

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type Task struct {
	Name      string
	JobId     string
	StartTime *time.Time
	StopTime  *time.Time
	ExitCode  int
}

type RunLog struct {
	RunId string
	State string
	Tasks []Task
}

func (m *Manager) GetWorkflowTasks(runId string) ([]Task, error) {
	m.readProjectSpec()
	m.readConfig()
	m.setContextForRun(runId)
	m.setContext(m.runContextName)
	m.setEngineForWorkflowType(m.runContextName)
	m.setContextStackInfo(m.runContextName)
	m.setWesUrl()
	m.setWesClient()
	m.getRunLog(runId)

	return m.getTasks()

}

func (m *Manager) GetRunLog(runId string) (RunLog, error) {
	if m.err != nil {
		return RunLog{}, m.err
	}
	var tasks []Task
	tasks, m.err = m.GetWorkflowTasks(runId)
	if m.err != nil {
		return RunLog{}, m.err
	}

	return RunLog{
		RunId: m.taskProps.runLog.RunId,
		State: string(m.taskProps.runLog.State),
		Tasks: tasks,
	}, nil
}

func (m *Manager) setContextForRun(runId string) {
	if m.err != nil {
		return
	}
	instance, err := m.Ddb.GetWorkflowInstanceById(context.Background(), m.projectSpec.Name, m.userId, runId)
	if err != nil {
		m.err = err
		return
	}
	m.runContextName = instance.ContextName
}

func (m *Manager) getRunLog(runId string) {
	if m.err != nil {
		return
	}
	m.runLog, m.err = m.wes.GetRunLog(context.Background(), runId)
}

func (m *Manager) getTasks() ([]Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	tasks := make([]Task, len(m.runLog.TaskLogs))
	for i, taskLog := range m.runLog.TaskLogs {
		taskName := taskLog.Name
		nameParts := strings.Split(taskName, "|")
		if len(nameParts) != 2 {
			return nil, fmt.Errorf("unable to parse job ID from task name '%s'", taskName)
		}
		tasks[i] = Task{
			Name:      nameParts[0],
			JobId:     nameParts[1],
			StartTime: parseLogTime(taskLog.StartTime),
			StopTime:  parseLogTime(taskLog.EndTime),
			ExitCode:  int(taskLog.ExitCode),
		}
	}

	return tasks, nil
}

func parseLogTime(logTime string) *time.Time {
	if logTime == "" {
		return nil
	}
	isoTime, err := time.Parse(time.RFC3339, logTime)
	if err != nil {
		log.Debug().Msgf("Unable to parse log time '%s' to ISO 8601, skipping", logTime)
	}
	return &isoTime
}
