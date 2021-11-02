package cfn

import "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

var ActiveStacksFilter = []types.StackStatus{
	types.StackStatusCreateInProgress,
	types.StackStatusCreateComplete,
	types.StackStatusRollbackInProgress,
	types.StackStatusRollbackFailed,
	types.StackStatusRollbackComplete,
	types.StackStatusDeleteInProgress,
	types.StackStatusDeleteFailed,
	types.StackStatusUpdateInProgress,
	types.StackStatusUpdateCompleteCleanupInProgress,
	types.StackStatusUpdateComplete,
	types.StackStatusUpdateRollbackInProgress,
	types.StackStatusUpdateRollbackFailed,
	types.StackStatusUpdateRollbackCompleteCleanupInProgress,
	types.StackStatusUpdateRollbackComplete,
	types.StackStatusReviewInProgress,
	types.StackStatusImportInProgress,
	types.StackStatusImportComplete,
	types.StackStatusImportRollbackInProgress,
	types.StackStatusImportRollbackFailed,
	types.StackStatusImportRollbackComplete,
}

var QueryableStacksMap = map[types.StackStatus]bool{
	types.StackStatusCreateComplete:                          true,
	types.StackStatusUpdateInProgress:                        true,
	types.StackStatusUpdateCompleteCleanupInProgress:         true,
	types.StackStatusUpdateComplete:                          true,
	types.StackStatusUpdateRollbackInProgress:                true,
	types.StackStatusUpdateRollbackFailed:                    true,
	types.StackStatusUpdateRollbackCompleteCleanupInProgress: true,
	types.StackStatusUpdateRollbackComplete:                  true,
	types.StackStatusReviewInProgress:                        true,
	types.StackStatusImportInProgress:                        true,
	types.StackStatusImportComplete:                          true,
	types.StackStatusImportRollbackInProgress:                true,
	types.StackStatusImportRollbackFailed:                    true,
	types.StackStatusImportRollbackComplete:                  true,
}
