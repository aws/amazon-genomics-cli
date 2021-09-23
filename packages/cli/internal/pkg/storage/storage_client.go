// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type StorageInstance struct {
	fsUtils *afero.Afero
}

func NewStorageInstance(fsOptional ...afero.Fs) (*StorageInstance, error) {
	var fs afero.Fs
	if len(fsOptional) > 0 {
		fs = fsOptional[0]
	} else {
		fs = afero.NewOsFs()
	}
	fsUtils := &afero.Afero{Fs: fs}
	return &StorageInstance{fsUtils}, nil
}

// ReadAsBytes reads the specified file and returns its content as an array of bytes.
// The filename can be a URL that is of the form file://<absolute-file-path> or
// simply a filename.
func (si *StorageInstance) ReadAsBytes(url string) ([]byte, error) {
	data, err := si.fsUtils.ReadFile(stripFileURLPrefix(url))
	if err != nil {
		return nil, fmt.Errorf("couldn't read file %s: %w", url, err)
	}
	return data, nil
}

// ReadAsString reads the specified file and returns its content as a string.
// The filename can be a URL that is of the form file://<absolute-file-path> or
// simply a filename.
func (si *StorageInstance) ReadAsString(url string) (string, error) {
	data, err := si.ReadAsBytes(stripFileURLPrefix(url))
	return string(data), err
}

// WriteFromBytes writes an array of bytes to the specified file.
// The filename can be a URL that is of the form file://<absolute-file-path> or
// simply a filename. The file is created if it doesn't exist and is overwritten
// if it does already exist. The directory that the file is in is created if it doesn't
// already exist.
func (si *StorageInstance) WriteFromBytes(url string, data []byte) error {
	filename := stripFileURLPrefix(url)
	if err := si.fsUtils.MkdirAll(filepath.Dir(filename), 0755 /* -rwxr-xr-x */); err != nil {
		return fmt.Errorf("couldn't create directories for file %s: %w", filename, err)
	}
	if err := si.fsUtils.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("couldn't write file %s: %w", filename, err)
	}
	return nil
}

// WriteFromString writes a string to the specified file.
// The filename can be a URL that is of the form file://<absolute-file-path> or
// simply a filename. The file is created if it doesn't exist and is overwritten
// if it does already exist. The directory that the file is in is created if it doesn't
// already exist.
func (si *StorageInstance) WriteFromString(url string, data string) error {
	return si.WriteFromBytes(url, []byte(data))
}

func stripFileURLPrefix(filename string) string {
	if strings.HasPrefix(filename, "file://") {
		runes := []rune(filename)
		filename = string(runes[7:])
	}
	return filename
}
