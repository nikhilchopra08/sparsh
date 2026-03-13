package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-production-secret-key-change-me")

type Claims struct {
	LicenseID int64  `json:"license_id"`
	DeviceID  string `json:"device_id"`
	jwt.RegisteredClaims
}

func GenerateToken(licenseID int64, deviceID string) (string, error) {
	expirationTime := time.Now().Add(24 * 365 * 10 * time.Hour) // 10 years for desktop app
	claims := &Claims{
		LicenseID: licenseID,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
