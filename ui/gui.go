// File: ui/gui.go

package ui

import (
	"log"
	"os"

	"govcs/config"
	"govcs/sync"

	"fyne.io/fyne"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func StartApp() {
	a := app.New()
	w := a.NewWindow("GoVCS - Google Drive VCS Tool")

	// Ensure config directory exists
	err := os.MkdirAll("config", os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	// Load existing config if available
	err = config.LoadConfig()
	if err != nil {
		log.Println("No existing config found. Starting fresh.")
	}

	localDirLabel := widget.NewLabel("Local Directory: " + config.AppConfig.LocalDir)
	remoteDirLabel := widget.NewLabel("Remote Directory ID: " + config.AppConfig.RemoteDir)

	selectLocalBtn := widget.NewButton("Select Local Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			localDir := uri.Path()
			config.AppConfig.LocalDir = localDir
			localDirLabel.SetText("Local Directory: " + localDir)
			err = config.SaveConfig()
			if err != nil {
				dialog.ShowError(err, w)
			}
		}, w)
	})

	selectRemoteBtn := widget.NewButton("Set Remote Directory ID", func() {
		dialog.ShowEntryDialog("Set Remote Directory ID", "Enter Remote Directory ID:", func(remoteDir string) {
			if remoteDir == "" {
				return
			}
			config.AppConfig.RemoteDir = remoteDir
			remoteDirLabel.SetText("Remote Directory ID: " + remoteDir)
			err := config.SaveConfig()
			if err != nil {
				dialog.ShowError(err, w)
			}
		}, w)
	})

	statusLabel := widget.NewLabel("Status: Ready")

	syncBtn := widget.NewButton("Sync Now", func() {
		if config.AppConfig.LocalDir == "" || config.AppConfig.RemoteDir == "" {
			dialog.ShowInformation("Error", "Please select both local and remote directories.", w)
			return
		}

		statusLabel.SetText("Status: Syncing...")
		go func() {
			sync.SyncDirectories()
			statusLabel.SetText("Status: Sync Completed")
		}()
	})

	content := container.NewVBox(
		localDirLabel,
		selectLocalBtn,
		remoteDirLabel,
		selectRemoteBtn,
		syncBtn,
		statusLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(500, 250))
	w.ShowAndRun()
}
