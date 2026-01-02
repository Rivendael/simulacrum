package data

import (
	"fmt"
)

type DriverLicense struct {
	Number         string `json:"number"`
	IssueDate      string `json:"issue_date"`
	ExpirationDate string `json:"expiration_date"`
}

// GenerateDeterministicDriverLicenseNumber generates a deterministic driver's license number
// Format: 3 letters + 6 digits (e.g., ABC123456)
func GenerateDeterministicDriverLicenseNumber(id, realDriverLicenseNumber string) string {
	if realDriverLicenseNumber == "" {
		return ""
	}
	hash := hashField(id, "driverlicense", realDriverLicenseNumber)

	// First three characters (letters A-Z)
	letter1 := 'A' + rune(hash[0]%26)
	letter2 := 'A' + rune(hash[1]%26)
	letter3 := 'A' + rune(hash[2]%26)

	// Six digits
	digits := (int(hash[3])<<16 | int(hash[4])<<8 | int(hash[5])) % 1000000

	return fmt.Sprintf("%c%c%c%06d", letter1, letter2, letter3, digits)
}

// ObscureDriverLicense takes real driver license data and returns a deterministic fake license
func ObscureDriverLicense(id string, real *DriverLicense) *DriverLicense {
	if real == nil {
		return nil
	}

	return &DriverLicense{
		Number:         GenerateDeterministicDriverLicenseNumber(id, real.Number),
		IssueDate:      GenerateDeterministicDate(id, "license_issue", real.IssueDate),
		ExpirationDate: GenerateDeterministicDate(id, "license_expiry", real.ExpirationDate),
	}
}
