package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestProjectValidate_Execute(t *testing.T) {
	testCases := map[string]struct {
		expectedErr error
		createMocks func(opts *validateProjectOpts)
	}{
		"project valid": {
			expectedErr: nil,
			createMocks: func(opts *validateProjectOpts) {
				opts.projectClient.(*storagemocks.MockProjectClient).EXPECT().Read().Return(spec.Project{}, nil)
			},
		},
		"invalid project": {
			expectedErr: fmt.Errorf("invalid spec"),
			createMocks: func(opts *validateProjectOpts) {
				opts.projectClient.(*storagemocks.MockProjectClient).EXPECT().Read().Return(spec.Project{}, fmt.Errorf("invalid spec"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProjClient := storagemocks.NewMockProjectClient(ctrl)
			opts := &validateProjectOpts{projectClient: mockProjClient}

			tc.createMocks(opts)
			err := opts.Execute()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
