// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
)

// StorageClient is the interface to storage systems
// from which to obtain various definition files, etc.
type StorageClient interface {
	ReadAsBytes(url string) ([]byte, error)
	ReadAsString(url string) (string, error)
	WriteFromBytes(url string, data []byte) error
	WriteFromString(url string, data string) error
}

type ProjectClient interface {
	Read() (spec.Project, error)
	Write(projectSpec spec.Project) error
	IsInitialized() (bool, error)
	GetProjectName() (string, error)
	GetLocation() string
}
