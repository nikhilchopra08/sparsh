package license

import (
	"crypto/rand"
	"fmt"
)

// GenerateKey creates a cryptographically random license key in the format XXXX-XXXX-XXXX-XXXX
func GenerateKey() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Omitted confusing chars like O, 0, I, 1
	keyParts := make([]string, 4)
	
	for i := 0; i < 4; i++ {
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		
		part := ""
		for _, b := range bytes {
			part += string(charset[int(b)%len(charset)])
		}
		keyParts[i] = part
	}
	
	return fmt.Sprintf("%s-%s-%s-%s", keyParts[0], keyParts[1], keyParts[2], keyParts[3]), nil
}
