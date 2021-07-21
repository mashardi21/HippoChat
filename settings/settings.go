// Package settings is a package that handles the saving and loading of a settings.json
// file for the HippoChat program.
package settings

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var (
	configPath, _ = os.UserConfigDir()
	settingsPath  = filepath.Join(configPath, "HippoChat", "settings.json")
	isConfig      bool
)

type Settings struct {
	Username string
	ID       string
}

func init() {
	isConfig = checkForConfig()
	if !isConfig {
		createDefaultConfig()
	}
}

// checkForConfig check that status of the settings.json file
// to determine if the file exists
func checkForConfig() bool {
	_, err := os.Stat(settingsPath)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Fatalf("There was an issue accesing the settings file: %s", err)
	}

	return true
}

// createDefaultConfig creates a Settings object containing
// default settings, marshals those settings into JSON, and
// then creates a new settings.json file containing the
// default settings
func createDefaultConfig() {
	defaults := &Settings{Username: "Egbog", ID: "Default"}

	data, err := json.Marshal(defaults)
	if err != nil {
		log.Fatalf("Could not encode settings into json: %s", err)
	}

	if _, err := os.Stat(filepath.Join(configPath, "HippoChat")); os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(configPath, "HippoChat"), 0777)
		if err != nil {
			log.Fatalf("Could not create directory: %s", err)
		}

	}

	err = os.WriteFile(settingsPath, data, 0664)
	if err != nil {
		log.Printf("Cannot write to file: %s", err)
	}

}

// Save takes a Settings object and saves the settings
// to settings.json
func Save(settings *Settings) {
	data, err := json.Marshal(settings)
	if err != nil {
		log.Printf("Could not convert settings to JSON: %s", err)
	}
	err = os.WriteFile(settingsPath, data, 0664)
	if err != nil {
		log.Printf("Cannot write to save file: %s", err)
	}
}

// Load reads the settings.json file and unmarshals the
// JSON into a Settings object that gets returned.
func Load() (*Settings, error) {
	var settings *Settings

	file, err := os.ReadFile(settingsPath)
	if err != nil {
		log.Printf("Could not access save file: %s", err)
		return nil, err
	}

	err = json.Unmarshal(file, &settings)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		return nil, err
	}

	return settings, nil
}
