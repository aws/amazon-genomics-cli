// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"reflect"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/stretchr/testify/require"
)

func TestConfigureDescribeContextOpts_Execute(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	describeContextOpts, err := newConfigureDescribeContextOpts()
	require.NoError(t, err)

	mockedConfig := config.Config{
		User: config.User{
			Email: "test-email@amazon.com",
			Id:    "testemail123",
		},
		Format: config.Format{
			Value: "text",
		},
	}

	mocks.configMock.EXPECT().Read().Return(mockedConfig, nil)
	describeContextOpts.configClient = mocks.configMock
	configReturned, err := describeContextOpts.Execute()
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(configReturned, mockedConfig))
}
