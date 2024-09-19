// File: sync/upload.go

package sync

import (
	"log"

	"govcs/drive"
)

// UploadFiles handles uploading a list of files to Google Drive.
func UploadFiles(files []string) {
	for _, file := range files {
		err := drive.UploadFile(file)
		if err != nil {
			log.Printf("Failed to upload %s: %v", file, err)
		} else {
			log.Printf("Successfully uploaded %s", file)
		}
	}
}
