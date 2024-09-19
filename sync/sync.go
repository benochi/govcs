// File: sync/sync.go

package sync

import (
	"log"

	"govcs/config"
	"govcs/drive"
)

// SyncDirectories compares local and remote directories and syncs changes.
func SyncDirectories() {
	// Authenticate with Google Drive
	err := drive.Authenticate()
	if err != nil {
		log.Fatalf("Drive authentication failed: %v", err)
	}

	// List remote files
	remoteFiles, err := drive.ListFiles()
	if err != nil {
		log.Fatalf("Failed to list remote files: %v", err)
	}

	// Compare files to determine which need to be uploaded
	filesToUpload := CompareFiles(config.AppConfig.LocalDir, remoteFiles)
	if len(filesToUpload) > 0 {
		log.Println("Starting upload of changed/new files...")
		UploadFiles(filesToUpload)
	} else {
		log.Println("No files to upload.")
	}

	// Optionally, implement pull logic here
	// For example, identify remote files not present locally or updated remotely
	// and download them using DownloadFiles
}
