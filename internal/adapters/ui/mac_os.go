package ui

import (
	"encoding/json"
	"file-mod-tracker/internal/adapters/config"
	"file-mod-tracker/internal/ports"
	"fmt"
	"fyne.io/fyne/v2"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MacOSUI struct {
	fileMonitorService ports.FileMonitorService
	workerAdapter      ports.WorkerAdapter
}

func NewMacOSUI(fileMonitorService ports.FileMonitorService, workerAdapter ports.WorkerAdapter) *MacOSUI {
	return &MacOSUI{
		fileMonitorService: fileMonitorService,
		workerAdapter:      workerAdapter,
	}
}

func (ui *MacOSUI) Show() {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Modification Tracker")

	messageLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	logsArea := widget.NewMultiLineEntry()
	logsArea.SetPlaceHolder("Logs will be displayed here...")
	logsArea.SetMinRowsVisible(10)

	scrollContainer := container.NewScroll(logsArea)
	scrollContainer.SetMinSize(fyne.NewSize(800, 400))

	statusLabel := widget.NewLabelWithStyle("Service Status: Stopped", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	showMessage := func(msg string, isError bool) {
		messageLabel.SetText(msg)
		messageLabel.TextStyle = fyne.TextStyle{Bold: true}
		if isError {
			messageLabel.Importance = widget.HighImportance
		} else {
			messageLabel.Importance = widget.MediumImportance
		}
	}

	startBtn := widget.NewButtonWithIcon("Start Service", theme.MediaPlayIcon(), func() {
		_, err := config.LoadConfig()
		if err != nil {
			showMessage(fmt.Sprintf("Error loading config: %v", err), true)
			return
		}
		// Implement start service logic here
		showMessage("Service started successfully!", false)
		statusLabel.SetText("Service Status: Running")
	})

	stopBtn := widget.NewButtonWithIcon("Stop Service", theme.MediaStopIcon(), func() {
		// Implement stop service logic here
		showMessage("Service stopped successfully!", false)
		statusLabel.SetText("Service Status: Stopped")
	})

	logsBtn := widget.NewButtonWithIcon("Fetch Logs", theme.DocumentIcon(), func() {
		fileChanges := ui.workerAdapter.GetFileChanges()
		if len(fileChanges) == 0 {
			logsArea.SetText("No file changes detected yet.")
		} else {
			logsJSON, err := json.MarshalIndent(fileChanges, "", "  ")
			if err != nil {
				log.Println("Error converting logs to JSON:", err)
				showMessage(fmt.Sprintf("Error converting logs to JSON: %v", err), true)
				return
			}
			logsArea.SetText(string(logsJSON))
		}
		showMessage("Logs fetched successfully!", false)
	})

	buttonContainer := container.NewHBox(startBtn, stopBtn, logsBtn)

	content := container.NewVBox(
		statusLabel,
		buttonContainer,
		messageLabel,
		scrollContainer,
	)
	content.Add(widget.NewSeparator())

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
