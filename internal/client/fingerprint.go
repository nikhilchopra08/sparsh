package client

import (
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"runtime"
)

// GetDeviceFingerprint generates a unique ID for the device based on hardware
func GetDeviceFingerprint() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	mac := getMacAddress()
	cpu := runtime.GOARCH + "_" + runtime.GOOS // Simplified CPU info if gopsutil is too heavy

	data := fmt.Sprintf("%s|%s|%s", hostname, mac, cpu)
	hash := sha256.Sum256([]byte(data))

	return fmt.Sprintf("%x", hash), nil
}

func getMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "00:00:00:00:00:00"
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && i.HardwareAddr.String() != "" {
			return i.HardwareAddr.String()
		}
	}
	return "00:00:00:00:00:00"
}
