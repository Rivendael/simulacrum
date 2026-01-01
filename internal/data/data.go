package data

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type PersonalData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	TaxID       string `json:"tax_id"`
}

// ObscureData takes the real data and returns a deterministic false identity
func ObscureData(real PersonalData) PersonalData {
	fake := PersonalData{
		ID: real.ID,
	}

	fake.Name = GenerateDeterministicName(real.ID, real.Name)
	fake.Email = GenerateDeterministicEmail(real.ID, real.Email)
	fake.Address = GenerateDeterministicAddress(real.ID, real.Address)
	fake.PhoneNumber = GenerateDeterministicPhone(real.ID, real.PhoneNumber)
	fake.TaxID = GenerateDeterministicTaxID(real.ID, real.TaxID)

	return fake
}

func GenerateDeterministicName(id, realName string) string {
	if realName == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(id + ":name:" + realName))

	firstIdx := int(hash[0]) % len(FirstNames)
	lastIdx := int(hash[1]) % len(LastNames)

	return fmt.Sprintf("%s %s", FirstNames[firstIdx], LastNames[lastIdx])
}

func GenerateDeterministicEmail(id, realEmail string) string {
	if realEmail == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(id + ":email:" + realEmail))

	firstIdx := int(hash[0]) % len(FirstNames)
	lastIdx := int(hash[1]) % len(LastNames)
	domainsIdx := int(hash[2]) % 4

	first := strings.ToLower(FirstNames[firstIdx])
	last := strings.ToLower(LastNames[lastIdx])
	domains := []string{"example.com", "test.org", "fake.net", "mail.com"}

	return fmt.Sprintf("%s.%s@%s", first, last, domains[domainsIdx])
}

func GenerateDeterministicAddress(id, realAddress string) string {
	if realAddress == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(id + ":address:" + realAddress))

	number := (int(hash[0])<<8 | int(hash[1])) % 9999
	if number == 0 {
		number = 1
	}

	streets := []string{"Main St", "Oak Ave", "Pine Ln", "Maple Dr", "Cedar Rd", "Elm St", "Washington Blvd"}
	streetIdx := int(hash[2]) % len(streets)

	return fmt.Sprintf("%d %s", number, streets[streetIdx])
}

func GenerateDeterministicPhone(id, realPhone string) string {
	if realPhone == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(id + ":phone:" + realPhone))

	exchange := int(hash[1])%900 + 100
	subscriber := (int(hash[2])<<8 | int(hash[3])) % 10000

	return fmt.Sprintf("555-%03d-%04d", exchange, subscriber)
}

func GenerateDeterministicTaxID(id, realTaxID string) string {
	if realTaxID == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(id + ":taxid:" + realTaxID))

	area := (int(hash[0])<<8 | int(hash[1])) % 900
	if area == 0 {
		area = 1
	}
	group := int(hash[2])%100 + 1
	serial := (int(hash[3])<<8 | int(hash[4])) % 10000

	return fmt.Sprintf("%03d-%02d-%04d", area, group, serial)
}
