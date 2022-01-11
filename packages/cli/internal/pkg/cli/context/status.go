package context

import "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
type Status string
const (
	StatusStarted    Status = "STARTED"
	StatusStopped    Status = "STOPPED"
	StatusFailed     Status = "FAILED"
	StatusNotStarted Status = "NOT_STARTED"
	StatusUnknown    Status = "UNKNOWN"
)

func (s Status) IsStarted() bool {
	return s == StatusStarted
}

func (s Status) IsStopped() bool {
	return s == StatusStopped
}

func (s Status) IsUnstarted() bool {
	return s == StatusNotStarted
}

func (s Status) IsFailed() bool {
	return s == StatusFailed
}

func (s Status) ToString() string {
	return string(s)
}

func mapStackToStatus(status types.StackStatus) Status {
	switch status {
	case "":
		return StatusNotStarted
	case types.StackStatusCreateInProgress,
		types.StackStatusCreateComplete,
		types.StackStatusUpdateInProgress,
		types.StackStatusUpdateCompleteCleanupInProgress,
		types.StackStatusUpdateComplete,
		types.StackStatusReviewInProgress,
		types.StackStatusImportInProgress,
		types.StackStatusImportComplete:
		return StatusStarted
	case types.StackStatusDeleteInProgress,
		types.StackStatusDeleteComplete:
		return StatusStopped
	case types.StackStatusCreateFailed,
		types.StackStatusRollbackInProgress,
		types.StackStatusRollbackFailed,
		types.StackStatusRollbackComplete,
		types.StackStatusDeleteFailed,
		types.StackStatusUpdateRollbackInProgress,
		types.StackStatusUpdateRollbackFailed,
		types.StackStatusUpdateRollbackCompleteCleanupInProgress,
		types.StackStatusUpdateRollbackComplete,
		types.StackStatusImportRollbackInProgress,
		types.StackStatusImportRollbackFailed,
		types.StackStatusImportRollbackComplete:
		return StatusFailed
	default:
		return StatusUnknown
	}
}
