// File: sync/compare.go

package sync

import (
	"log"
	"os"
	"path/filepath"
	"time"

	driveapi "google.golang.org/api/drive/v3"
)

// CompareFiles determines which local files need to be uploaded based on modification times.
func CompareFiles(localDir string, remoteFiles []*driveapi.File) []string {
	var filesToUpload []string

	// Create a map for quick lookup of remote files by name
	remoteMap := make(map[string]*driveapi.File)
	for _, file := range remoteFiles {
		remoteMap[file.Name] = file
	}

	// Walk through the local directory
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		// Skip directories
		if !info.IsDir() {
			// Get the relative path of the file
			relPath, err := filepath.Rel(localDir, path)
			if err != nil {
				log.Printf("Failed to get relative path for %s: %v", path, err)
				return err
			}

			remoteFile, exists := remoteMap[relPath]

			if !exists {
				// File doesn't exist remotely; needs to be uploaded
				filesToUpload = append(filesToUpload, path)
			} else {
				// Parse the remote file's ModifiedTime
				remoteModTime, err := time.Parse(time.RFC3339, remoteFile.ModifiedTime)
				if err != nil {
					log.Printf("Failed to parse ModifiedTime for %s: %v", remoteFile.Name, err)
					return err
				}

				// Compare modification times
				if info.ModTime().After(remoteModTime) {
					filesToUpload = append(filesToUpload, path)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %q: %v\n", localDir, err)
	}

	return filesToUpload
}
