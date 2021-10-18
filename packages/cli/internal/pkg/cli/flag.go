// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import "github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"

const (
	argsFlag            = "args"
	argsFlagShort       = "a"
	argsFlagDescription = "Arguments to use."
)

const (
	VerboseFlag            = "verbose"
	VerboseFlagShort       = "v"
	VerboseFlagDescription = "Display verbose diagnostic information."
)

const (
	FormatFlag            = "format"
	FormatFlagDefault     = string(format.DefaultFormat)
	FormatFlagDescription = "Format option for output. Valid options are: text, tabular"
)

const (
	AWSProfileFlag            = "awsProfile"
	AWSProfileFlagShort       = "p"
	AWSProfileFlagDescription = "Use the provided AWS CLI profile."
)
