package drive

import (
	"io"
	"os"
	"path/filepath"

	"govcs/config"

	driveapi "google.golang.org/api/drive/v3"
)

// ListFiles retrieves a list of files in the specified remote directory.
func ListFiles() ([]*driveapi.File, error) {
	query := "'" + config.AppConfig.RemoteDir + "' in parents and trashed=false"
	req := Service.Files.List().Q(query).Fields("files(id, name, modifiedTime)")
	res, err := req.Do()
	if err != nil {
		return nil, err
	}
	return res.Files, nil
}

// UploadFile uploads a single file to Google Drive.
func UploadFile(localPath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fileName := filepath.Base(localPath)
	file := &driveapi.File{
		Name:    fileName,
		Parents: []string{config.AppConfig.RemoteDir},
	}

	_, err = Service.Files.Create(file).Media(f).Do()
	return err
}

// DownloadFile downloads a single file from Google Drive.
func DownloadFile(fileID, localPath string) error {
	resp, err := Service.Files.Get(fileID).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
