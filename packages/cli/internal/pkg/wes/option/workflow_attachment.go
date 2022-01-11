package option

import (
	"os"
	"github.com/antihax/optional"
	wes "github.com/rsc/wes_client"
)

func WorkflowAttachment(attachmentPaths []string) Func {
	return func(opts *wes.RunWorkflowOpts) error {
		var fileDescriptors []*os.File
		for _, attachmentPath := range attachmentPaths {
			fileDescriptor, err := os.Open(attachmentPath)
			if err != nil {
				return err
			}
			fileDescriptors = append(fileDescriptors, fileDescriptor)
		}
		opts.WorkflowAttachment = optional.NewInterface(fileDescriptors)
		return nil
	}
}
