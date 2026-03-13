package models

import "time"

type License struct {
	ID         int64     `json:"id" db:"id"`
	LicenseKey string    `json:"license_key" db:"license_key"`
	Email      string    `json:"email" db:"email"`
	Status     string    `json:"status" db:"status"` // "active", "revoked", "expired"
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type Activation struct {
	ID          int64     `json:"id" db:"id"`
	LicenseID   int64     `json:"license_id" db:"license_id"`
	DeviceID    string    `json:"device_id" db:"device_id"`
	ActivatedAt time.Time `json:"activated_at" db:"activated_at"`
}

type ActivateRequest struct {
	LicenseKey string `json:"license_key" binding:"required"`
	DeviceID   string `json:"device_id" binding:"required"`
}

type ActivateResponse struct {
	ActivationToken string `json:"activation_token"`
}

type UserSettings struct {
	CompanyName    string `json:"company_name"`
	CompanyAddress string `json:"company_address"`
	LogoPath       string `json:"logo_path"`
	LastInvoiceNum string `json:"last_invoice_num"`
}

type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Rate        float64 `json:"rate"`
	Amount      float64 `json:"amount"`
}

type Invoice struct {
	Number       string        `json:"number"`
	Date         string        `json:"date"`
	DueDate      string        `json:"due_date"`
	PaymentTerms string        `json:"payment_terms"`
	PONumber     string        `json:"po_number"`
	From         string        `json:"from"`
	BillToName   string        `json:"bill_to_name"`
	BillToAddr   string        `json:"bill_to_addr"`
	ShipToName   string        `json:"ship_to_name"`
	ShipToAddr   string        `json:"ship_to_addr"`
	Items        []InvoiceItem `json:"items"`
	Notes        string        `json:"notes"`
	Terms        string        `json:"terms"`
	TaxRate      float64       `json:"tax_rate"`
	AmountPaid   float64       `json:"amount_paid"`
	LogoPath     string        `json:"logo_path"`
}
