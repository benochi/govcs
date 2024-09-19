// File: sync/download.go

package sync

import (
	"log"
	"path/filepath"

	"govcs/config"
	"govcs/drive"
	"govcs/utils"

	driveapi "google.golang.org/api/drive/v3"
)

// DownloadFiles handles downloading a list of files from Google Drive.
func DownloadFiles(files []*driveapi.File) {
	backupDir := filepath.Join(config.AppConfig.LocalDir, "backup")
	err := utils.BackupFiles(config.AppConfig.LocalDir, backupDir)
	if err != nil {
		log.Printf("Failed to backup local files: %v", err)
		return
	}

	for _, file := range files {
		localPath := filepath.Join(config.AppConfig.LocalDir, file.Name)
		err := drive.DownloadFile(file.Id, localPath)
		if err != nil {
			log.Printf("Failed to download %s: %v", file.Name, err)
		} else {
			log.Printf("Successfully downloaded %s", file.Name)
		}
	}
}
