// File: sync/compare.go

package sync

import (
	"log"
	"os"
	"path/filepath"

	driveapi "google.golang.org/api/drive/v3"
)

// CompareFiles determines which local files need to be uploaded based on modification times.
func CompareFiles(localDir string, remoteFiles []*driveapi.File) []string {
	var filesToUpload []string

	remoteMap := make(map[string]*driveapi.File)
	for _, file := range remoteFiles {
		remoteMap[file.Name] = file
	}

	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(localDir, path)
			if err != nil {
				log.Printf("Failed to get relative path for %s: %v", path, err)
				return err
			}
			remoteFile, exists := remoteMap[relPath]

			if !exists || info.ModTime().After(remoteFile.ModifiedTime.AsTime()) {
				filesToUpload = append(filesToUpload, path)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %q: %v\n", localDir, err)
	}

	return filesToUpload
}
