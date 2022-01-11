// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"testing"
	"github.com/stretchr/testify/require"
)

const (
	correctEmailAddress   = "my@email.com"
	incorrectEmailAddress = "invalid email address"
)

func TestEmailContextOpts_Execute(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	emailContextOpts, err := newEmailContextOpts(emailContextVars{
		userEmailAddress: correctEmailAddress,
	})
	require.NoError(t, err)

	mocks.configMock.EXPECT().SetUserEmailAddress(correctEmailAddress).Return(nil)
	emailContextOpts.configClient = mocks.configMock
	err = emailContextOpts.Execute()

	require.NoError(t, err)
}

func TestEmailContextOpts_Validate_CorrectAddress(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	emailContextOpts, err := newEmailContextOpts(emailContextVars{
		userEmailAddress: correctEmailAddress,
	})
	require.NoError(t, err)

	emailContextOpts.configClient = mocks.configMock
	err = emailContextOpts.Validate()

	require.NoError(t, err)
}

func TestEmailContextOpts_Validate_IncorrectAddress(t *testing.T) {
	mocks := createMocks(t)
	defer mocks.ctrl.Finish()

	emailContextOpts, err := newEmailContextOpts(emailContextVars{
		userEmailAddress: incorrectEmailAddress,
	})
	require.NoError(t, err)

	emailContextOpts.configClient = mocks.configMock
	err = emailContextOpts.Validate()

	require.Error(t, err)
}
