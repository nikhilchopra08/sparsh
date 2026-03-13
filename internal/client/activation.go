package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"licensebox/internal/auth"
	"licensebox/internal/models"
)

const TokenFile = ".license_token"

type LicenseManager struct {
	ServerURL string
}

func NewLicenseManager(serverURL string) *LicenseManager {
	return &LicenseManager{ServerURL: serverURL}
}

// CheckLicense flows through the activation check
func (m *LicenseManager) CheckLicense() (bool, error) {
	tokenPath := filepath.Join(".", TokenFile)

	// 1. Check if token exists
	if _, err := os.Stat(tokenPath); err == nil {
		data, err := ioutil.ReadFile(tokenPath)
		if err == nil {
			// 2. Validate token signature
			claims, err := auth.ValidateToken(string(data))
			if err == nil {
				// Check if this token belongs to THIS device
				deviceID, _ := GetDeviceFingerprint()
				if claims.DeviceID == deviceID {
					fmt.Println("✅ License validated successfully (Offline)")
					return true, nil
				}
			}
		}
	}

	// 3. Prompt or return false
	return false, nil
}

// Activate requests a new token from the server
func (m *LicenseManager) Activate(licenseKey string) error {
	deviceID, err := GetDeviceFingerprint()
	if err != nil {
		return err
	}

	reqBody := models.ActivateRequest{
		LicenseKey: licenseKey,
		DeviceID:   deviceID,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(m.ServerURL+"/activate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("activation failed: %s", errResp["error"])
	}

	var res models.ActivateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	// 4. Store token locally
	return ioutil.WriteFile(TokenFile, []byte(res.ActivationToken), 0600)
}
