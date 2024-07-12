package main

import (
	"errors"
	"os"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var appFolder string
var logsFolder string
var savedConfigFolder string
var configPath string
var appIconPath string

func path_init() error {
	appData, err := os.UserConfigDir()
	if err != nil {
		appData = os.Getenv("APPDATA")
		if appData == "" {
			return errors.New("Could not find user config directory: " + err.Error())
		}
	}
	runtime.LogDebug(appContext, "Found user config directory: "+appData)

	appFolder = path.Join(appData, "folder-creator")
	logsFolder = path.Join(appFolder, "logs")
	savedConfigFolder = path.Join(appFolder, "savedconfigs")

	configPath = path.Join(appFolder, "config.json")
	appIconPath = path.Join(appFolder, "appicon.png")

	runtime.LogTrace(appContext, "Attempting to create folders")
	err = create_folder(appFolder)
	if err != nil {
		return err
	}

	err = create_folder(logsFolder)
	if err != nil {
		return err
	}
	err = create_folder(savedConfigFolder)
	if err != nil {
		return err
	}

	runtime.LogTrace(appContext, "Creating folders complete")

	runtime.LogTrace(appContext, "Attempting to create appicon")

	// Create icon from embedded appIcon if it exists
	if _, err := os.Stat(appIconPath); os.IsNotExist(err) {
		runtime.LogTrace(appContext, "appicon not found, creating from embedded appIcon")
		err = os.WriteFile(appIconPath, appIcon, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}

func get_logs_folder() (string, error) {
	logsFolder = path.Join(os.Getenv("APPDATA"), "folder-creator", "logs")

	// Create folder if it doesn't exist
	if _, err := os.Stat(logsFolder); os.IsNotExist(err) {
		err = os.MkdirAll(logsFolder, 0o755)
		if err != nil {
			return "", err
		}
	}
	return logsFolder, nil
}

func get_config_path() string {
	configPath = path.Join(os.Getenv("APPDATA"), "folder-creator", "config.json")

	return configPath
}

// Create folder if it doesn't exist, return error
func create_folder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0o755)
		if err != nil {
			return err
		}
	} else {
		runtime.LogDebug(appContext, "Folder already exists: "+folder)
		return nil
	}
	runtime.LogDebug(appContext, "Created folder: "+folder)

	return nil
}
