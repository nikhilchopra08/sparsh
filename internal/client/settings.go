package client

import (
	"encoding/json"
	"io/ioutil"
	"licensebox/internal/models"
	"os"
)

const SettingsFile = ".user_settings.json"

func SaveUserSettings(settings models.UserSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(SettingsFile, data, 0644)
}

func LoadUserSettings() (*models.UserSettings, error) {
	if _, err := os.Stat(SettingsFile); os.IsNotExist(err) {
		return &models.UserSettings{}, nil
	}

	data, err := ioutil.ReadFile(SettingsFile)
	if err != nil {
		return nil, err
	}

	var settings models.UserSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

func SaveInvoice(invoice models.Invoice, path string) error {
	data, err := json.MarshalIndent(invoice, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func LoadInvoice(path string) (*models.Invoice, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var inv models.Invoice
	if err := json.Unmarshal(data, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}
