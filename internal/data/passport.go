package data

import (
	"fmt"
)

type Passport struct {
	Number         string `json:"number"`
	IssueDate      string `json:"issue_date"`
	ExpirationDate string `json:"expiration_date"`
}

// GenerateDeterministicPassportNumber generates a deterministic passport number
// Format: 2 letters + 7 digits (e.g., AB1234567)
func GenerateDeterministicPassportNumber(id, realPassportNumber string) string {
	if realPassportNumber == "" {
		return ""
	}
	hash := hashField(id, "passport", realPassportNumber)

	// First two characters (letters A-Z)
	letter1 := 'A' + rune(hash[0]%26)
	letter2 := 'A' + rune(hash[1]%26)

	// Seven digits
	digits := (int(hash[2])<<24 | int(hash[3])<<16 | int(hash[4])<<8 | int(hash[5])) % 10000000

	return fmt.Sprintf("%c%c%07d", letter1, letter2, digits)
}

// ObscurePassport takes real passport data and returns a deterministic fake passport
func ObscurePassport(id string, real *Passport) *Passport {
	if real == nil {
		return nil
	}

	return &Passport{
		Number:         GenerateDeterministicPassportNumber(id, real.Number),
		IssueDate:      GenerateDeterministicDate(id, "passport_issue", real.IssueDate),
		ExpirationDate: GenerateDeterministicDate(id, "passport_expiry", real.ExpirationDate),
	}
}
