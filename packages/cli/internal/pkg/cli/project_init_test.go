// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectInit_Validate(t *testing.T) {
	testCases := map[string]struct {
		projectName  string
		expectedErr  string
		workflowType string
	}{
		"valid project name": {
			projectName:  testProjectName,
			workflowType: "nextflow",
		},
		"illegal project name": {
			projectName:  testBadProjectName,
			expectedErr:  fmt.Sprintf("%s has non-alpha-numeric characters in it", testBadProjectName),
			workflowType: "wdl",
		},
		"illegal project name 2": {
			projectName:  testBadProjectName2,
			expectedErr:  fmt.Sprintf("%s has non-alpha-numeric characters in it", testBadProjectName2),
			workflowType: "wdl",
		},
		"missing project name": {
			expectedErr:  "missing project name",
			workflowType: "nextflow",
		},
		"invalid workflow type": {
			expectedErr:  "invalid workflow type supplied: 'aBadEngineName'. Supported workflow types are: [nextflow snakemake wdl]",
			workflowType: "aBadEngineName",
			projectName:  testProjectName,
		},
		"missing workflow type": {
			expectedErr:  "please specify a workflow type with the --workflow-type flag",
			workflowType: "",
			projectName:  testProjectName,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProj := storagemocks.NewMockProjectClient(ctrl)
			opts := &initProjectOpts{
				projectClient:   mockProj,
				initProjectVars: initProjectVars{tc.projectName, tc.workflowType},
			}
			mockProj.EXPECT().IsInitialized().AnyTimes().Return(false, nil)
			err := opts.Validate()

			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestProjectInit_Execute(t *testing.T) {
	testCases := map[string]struct {
		projectName    string
		expectedEngine []spec.Engine
		engineName     string
	}{
		"cromwell engine generation": {
			projectName:    testProjectName,
			engineName:     "wdl",
			expectedEngine: []spec.Engine{{Type: "wdl", Engine: "cromwell"}},
		},
		"nextflow engine generation": {
			projectName:    testProjectName,
			engineName:     "nextflow",
			expectedEngine: []spec.Engine{{Type: "nextflow", Engine: "nextflow"}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProj := storagemocks.NewMockProjectClient(ctrl)

			expectedProject := spec.Project{
				Name:          tc.projectName,
				SchemaVersion: spec.LatestVersion,
				Contexts: map[string]spec.Context{
					initContextName: {RequestSpotInstances: false, Engines: tc.expectedEngine},
				},
			}

			opts := &initProjectOpts{
				projectClient:   mockProj,
				initProjectVars: initProjectVars{tc.projectName, tc.engineName},
			}
			mockProj.EXPECT().Write(expectedProject).AnyTimes().Return(nil)

			err := opts.Execute()

			require.NoError(t, err)
		})
	}
}
func TestProjectInit_CreateInitialProject_ValidateSchema(t *testing.T) {
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, storage.ProjectSpecFileName)
	specFile, err := os.Create(tempFilePath)
	if err != nil {
		t.Fatal(err)
	}
	_ = specFile.Close()
	client, err := storage.NewProjectClientWithLocation(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	opts := &initProjectOpts{
		initProjectVars: initProjectVars{
			ProjectName:  "AGCTestProject",
			workflowType: "wdl",
		},
		projectClient: client,
	}
	fileComponents := opts.createInitialProject()
	err = client.Write(fileComponents)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		t.Fatal(err)
	}
	err = spec.ValidateProject(bytes)
	require.NoError(t, err)
}
