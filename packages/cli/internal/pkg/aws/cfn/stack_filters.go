package cfn

import "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

type StackOptions struct {
	activeStack    bool
	queryableStack bool
}

var stackDefinitions = map[types.StackStatus]StackOptions{
	types.StackStatusCreateInProgress:                        {activeStack: true, queryableStack: false},
	types.StackStatusCreateComplete:                          {activeStack: true, queryableStack: true},
	types.StackStatusRollbackInProgress:                      {activeStack: true, queryableStack: false},
	types.StackStatusRollbackFailed:                          {activeStack: true, queryableStack: false},
	types.StackStatusRollbackComplete:                        {activeStack: true, queryableStack: false},
	types.StackStatusDeleteInProgress:                        {activeStack: true, queryableStack: false},
	types.StackStatusDeleteFailed:                            {activeStack: true, queryableStack: false},
	types.StackStatusUpdateInProgress:                        {activeStack: true, queryableStack: true},
	types.StackStatusUpdateCompleteCleanupInProgress:         {activeStack: true, queryableStack: true},
	types.StackStatusUpdateComplete:                          {activeStack: true, queryableStack: true},
	types.StackStatusUpdateRollbackInProgress:                {activeStack: true, queryableStack: true},
	types.StackStatusUpdateRollbackFailed:                    {activeStack: true, queryableStack: true},
	types.StackStatusUpdateRollbackCompleteCleanupInProgress: {activeStack: true, queryableStack: true},
	types.StackStatusUpdateRollbackComplete:                  {activeStack: true, queryableStack: false},
	types.StackStatusReviewInProgress:                        {activeStack: true, queryableStack: true},
	types.StackStatusImportInProgress:                        {activeStack: true, queryableStack: true},
	types.StackStatusImportComplete:                          {activeStack: true, queryableStack: true},
	types.StackStatusImportRollbackInProgress:                {activeStack: true, queryableStack: true},
	types.StackStatusImportRollbackFailed:                    {activeStack: true, queryableStack: true},
	types.StackStatusImportRollbackComplete:                  {activeStack: true, queryableStack: true},
}
var ActiveStacksFilter []types.StackStatus
var QueryableStacksMap map[types.StackStatus]bool

func init() {
	QueryableStacksMap = make(map[types.StackStatus]bool)
	for stackStatus, stackOptions := range stackDefinitions {
		QueryableStacksMap[stackStatus] = stackOptions.queryableStack
		if stackOptions.activeStack {
			ActiveStacksFilter = append(ActiveStacksFilter, stackStatus)
		}
	}
}
