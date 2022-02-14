package zipfile

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CompressToTmp(srcPath string) (string, error) {
	packFile, err := ioutil.TempFile("", "workflow_*")
	if err != nil {
		return "", err
	}
	defer packFile.Close()
	zipWriter := zip.NewWriter(packFile)
	defer zipWriter.Close()
	if err := writeToZipRecursive(zipWriter, srcPath); err != nil {
		return "", err
	}
	return packFile.Name(), nil
}

func writeToZipRecursive(writer *zip.Writer, rootPath string) error {
	// Expand home directory path
	return filepath.WalkDir(rootPath, func(currentPath string, dirEntry fs.DirEntry, err error) error {
		if dirEntry == nil {
			// There are several use cases when it can happen:
			// 1. provided path doesn't exist
			// 2. file or sub-directory got deleted after being listed by WalkDir
			return fmt.Errorf("file '%s' doesn't exist", currentPath)
		}
		if !dirEntry.IsDir() {
			return writeFileToZip(writer, rootPath, currentPath)
		}
		return nil
	})
}

func writeFileToZip(writer *zip.Writer, rootPath, currentPath string) error {
	srcFile, err := os.Open(currentPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	var relativePath string
	if rootPath == currentPath {
		relativePath = filepath.Base(currentPath)
	} else {
		relativePath, err = filepath.Rel(rootPath, currentPath)
		if err != nil {
			return err
		}
	}
	dstFile, err := writer.Create(relativePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstFile, srcFile)
	return err
}
