// File: utils/utils.go

package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destDir := filepath.Dir(dst)
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// BackupFiles copies all files from localDir to backupDir.
func BackupFiles(localDir, backupDir string) error {
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(localDir, path)
			if err != nil {
				return err
			}
			backupPath := filepath.Join(backupDir, relPath)
			return CopyFile(path, backupPath)
		}
		return nil
	})
	return err
}
