// go:build !windows

// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package logging

// Log message prefixes.
const (
	successPrefix = "✔" //nolint:deadcode,varcheck
	errorPrefix   = "✘ "
	warningPrefix = "⚠️ "
	infoPrefix    = "𝒊 "
	debugPrefix   = "↓ "
	fatalPrefix   = "☠️ "
	panicPrefix   = "!! "
	tracePrefix   = "🔍 "
)
