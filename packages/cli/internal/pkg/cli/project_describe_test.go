// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/types"
	storagemocks "github.com/aws/amazon-genomics-cli/cli/internal/pkg/mocks/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestProjectDescribe_Execute(t *testing.T) {
	testCases := map[string]struct {
		expectedProject types.Project
		expectedErr     error
		setupMocks      func(opts *describeProjectOpts)
	}{
		"describe project": {
			expectedProject: types.Project{
				Name: testProjectName,
				Data: []types.Data{{Location: testDataName1, ReadOnly: true}},
			},
			setupMocks: func(opts *describeProjectOpts) {
				opts.projectClient.(*storagemocks.MockProjectClient).EXPECT().Read().Return(spec.Project{Name: testProjectName, Data: []spec.Data{{Location: testDataName1, ReadOnly: true}}}, nil)
			},
		},
		"project client error case": {
			expectedErr: fmt.Errorf("cannot read"),
			setupMocks: func(opts *describeProjectOpts) {
				opts.projectClient.(*storagemocks.MockProjectClient).EXPECT().Read().Return(spec.Project{}, fmt.Errorf("cannot read"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProjectClient := storagemocks.NewMockProjectClient(ctrl)
			opts := &describeProjectOpts{
				projectClient:       mockProjectClient,
				describeProjectVars: describeProjectVars{},
			}

			tc.setupMocks(opts)
			projectDefn, err := opts.Execute()

			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expectedProject, projectDefn)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
