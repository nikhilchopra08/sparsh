package db

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite" // CGO-free sqlite driver
	"licensebox/internal/models"
	"time"
)

type Database struct {
	Conn *sql.DB
}

func InitDB(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS licenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			license_key TEXT UNIQUE NOT NULL,
			email TEXT NOT NULL,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS activations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			license_id INTEGER NOT NULL,
			device_id TEXT NOT NULL,
			activated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (license_id) REFERENCES licenses(id)
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	return &Database{Conn: db}, nil
}

func (d *Database) CreateLicense(email, key string) (*models.License, error) {
	res, err := d.Conn.Exec("INSERT INTO licenses (email, license_key) VALUES (?, ?)", email, key)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &models.License{
		ID:         id,
		LicenseKey: key,
		Email:      email,
		Status:     "active",
		CreatedAt:  time.Now(),
	}, nil
}

func (d *Database) GetLicenseByKey(key string) (*models.License, error) {
	row := d.Conn.QueryRow("SELECT id, license_key, email, status, created_at FROM licenses WHERE license_key = ?", key)
	var l models.License
	err := row.Scan(&l.ID, &l.LicenseKey, &l.Email, &l.Status, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (d *Database) CreateActivation(licenseID int64, deviceID string) error {
	_, err := d.Conn.Exec("INSERT INTO activations (license_id, device_id) VALUES (?, ?)", licenseID, deviceID)
	return err
}

func (d *Database) CheckActivation(licenseID int64, deviceID string) (bool, error) {
	var count int
	err := d.Conn.QueryRow("SELECT COUNT(*) FROM activations WHERE license_id = ? AND device_id = ?", licenseID, deviceID).Scan(&count)
	return count > 0, err
}
