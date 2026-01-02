package data

import (
	"fmt"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type PersonalData struct {
	ID            string         `json:"id"`
	Name          string         `json:"name,omitempty"`
	FirstName     string         `json:"first_name,omitempty"`
	LastName      string         `json:"last_name,omitempty"`
	MiddleName    string         `json:"middle_name,omitempty"`
	Email         string         `json:"email,omitempty"`
	Address       string         `json:"address,omitempty"`
	Street        string         `json:"street,omitempty"`
	City          string         `json:"city,omitempty"`
	State         string         `json:"state,omitempty"`
	ZipCode       string         `json:"zip_code,omitempty"`
	County        string         `json:"county,omitempty"`
	Country       string         `json:"country,omitempty"`
	PhoneNumber   string         `json:"phone_number,omitempty"`
	TaxID         string         `json:"tax_id,omitempty"`
	DateOfBirth   string         `json:"date_of_birth,omitempty"`
	Gender        string         `json:"gender,omitempty"`
	SSN           string         `json:"ssn,omitempty"`
	Passport      *Passport      `json:"passport,omitempty"`
	DriverLicense *DriverLicense `json:"driver_license,omitempty"`
	BankAccounts  []*BankAccount `json:"bank_accounts,omitempty"`
	// Generic integer and float fields
	IntegerValue int64   `json:"integer_value,omitempty"`
	FloatValue   float64 `json:"float_value,omitempty"`
	// Generic nested object (recursive)
	Object  *PersonalData   `json:"object,omitempty"`
	Objects []*PersonalData `json:"objects,omitempty"`
}

// ObscureData takes the real data and returns a deterministic false identity
func ObscureData(real PersonalData) PersonalData {
	fake := PersonalData{
		ID: real.ID,
	}

	fake.Name, fake.FirstName, fake.LastName, fake.MiddleName = ObscureName(real.ID, real)

	fake.Email = GenerateDeterministicEmail(real.ID, real.Email)

	fake.Address, fake.Street, fake.City, fake.State, fake.ZipCode, fake.County, fake.Country = ObscureAddress(real.ID, real)

	fake.PhoneNumber = GenerateDeterministicPhone(real.ID, real.PhoneNumber)
	fake.TaxID = GenerateDeterministicTaxID(real.ID, real.TaxID)
	fake.DateOfBirth = GenerateDeterministicDateOfBirth(real.ID, real.DateOfBirth)
	fake.Gender = GenerateDeterministicGender(real.ID, real.Gender)
	fake.SSN = GenerateDeterministicSSN(real.ID, real.SSN)
	fake.Passport = ObscurePassport(real.ID, real.Passport)
	fake.DriverLicense = ObscureDriverLicense(real.ID, real.DriverLicense)
	fake.BankAccounts = ObscureBankAccounts(real.ID, real.BankAccounts)

	// Handle generic integer and float fields
	if real.IntegerValue != 0 {
		fake.IntegerValue = GenerateDeterministicInteger(real.ID, real.IntegerValue)
	}
	if real.FloatValue != 0 {
		fake.FloatValue = GenerateDeterministicFloat(real.ID, real.FloatValue)
	}

	// Handle nested generic objects (recursive)
	if real.Object != nil {
		// Use the nested object's ID if available, otherwise use parent ID
		nestedID := real.Object.ID
		if nestedID == "" {
			nestedID = real.ID + "_object"
		}
		obscuredNested := ObscureData(*real.Object)
		fake.Object = &obscuredNested
	}

	if len(real.Objects) > 0 {
		fake.Objects = make([]*PersonalData, len(real.Objects))
		for i, obj := range real.Objects {
			// Use object's ID if available, otherwise generate from parent
			objID := obj.ID
			if objID == "" {
				objID = real.ID + "_object_" + fmt.Sprintf("%d", i)
			}
			obscuredObj := ObscureData(*obj)
			fake.Objects[i] = &obscuredObj
		}
	}

	return fake
}

// hashField generates an xxHash for a field, returning 8 bytes
func hashField(id, fieldType, value string) [8]byte {
	h := xxhash.Sum64([]byte(id + ":" + fieldType + ":" + value))
	var result [8]byte
	for i := range 8 {
		result[i] = byte(h >> (8 * i))
	}
	return result
}

// selectFromList uses hash bytes to deterministically select from a list
func selectFromList(hash [8]byte, byteIndex int, list []string) string {
	idx := int(hash[byteIndex]) % len(list)
	return list[idx]
}

// bytesToInt converts two hash bytes to an integer
func bytesToInt(hash [8]byte, byteIndex int, max int) int {
	val := (int(hash[byteIndex])<<8 | int(hash[(byteIndex+1)%8])) % max
	if val == 0 {
		val = 1
	}
	return val
}

func GenerateDeterministicEmail(id, realEmail string) string {
	if realEmail == "" {
		return ""
	}
	hash := hashField(id, "email", realEmail)
	first := strings.ToLower(selectFromList(hash, 0, FirstNames))
	last := strings.ToLower(selectFromList(hash, 1, LastNames))
	domains := []string{"example.com", "test.org", "fake.net", "mail.com"}
	domain := selectFromList(hash, 2, domains)
	return fmt.Sprintf("%s.%s@%s", first, last, domain)
}

func GenerateDeterministicPhone(id, realPhone string) string {
	if realPhone == "" {
		return ""
	}

	// Check if phone starts with a country code (e.g., "+1-...")
	var countryCode string
	if strings.HasPrefix(realPhone, "+") {
		// Extract country code and the rest
		parts := strings.SplitN(realPhone, "-", 2)
		if len(parts) == 2 {
			countryCode = parts[0] // e.g., "+1"
		}
	}

	hash := hashField(id, "phone", realPhone)
	exchange := int(hash[1])%900 + 100
	subscriber := (int(hash[2])<<8 | int(hash[3])) % 10000
	generatedPhone := fmt.Sprintf("555-%03d-%04d", exchange, subscriber)

	// If there was a country code, preserve it
	if countryCode != "" {
		return countryCode + "-" + generatedPhone
	}
	return generatedPhone
}

func GenerateDeterministicCountryCode(id, realCountryCode string) string {
	if realCountryCode == "" {
		return ""
	}
	hash := hashField(id, "country_code", realCountryCode)
	// Generate country code from 1-999
	code := int(hash[0])*256 + int(hash[1])
	countryCode := code%999 + 1
	return fmt.Sprintf("+%d", countryCode)
}

func GenerateDeterministicTaxID(id, realTaxID string) string {
	if realTaxID == "" {
		return ""
	}
	hash := hashField(id, "taxid", realTaxID)
	area := bytesToInt(hash, 0, 900)
	group := int(hash[2])%100 + 1
	serial := (int(hash[3])<<8 | int(hash[4])) % 10000
	return fmt.Sprintf("%03d-%02d-%04d", area, group, serial)
}

// GenerateDeterministicDate generates a deterministic date in YYYY-MM-DD format
// Creates realistic dates within a range (issue dates: past 10 years, expiration dates: future 5-10 years)
func GenerateDeterministicDate(id, fieldType, realDate string) string {
	if realDate == "" {
		return ""
	}
	hash := hashField(id, fieldType, realDate)

	// For issue dates, use a date in the past (0-10 years ago)
	// For expiration dates, use a date in the future (5-10 years ahead)
	var baseYear, yearRange int
	if fieldType == "passport_issue" || fieldType == "license_issue" {
		baseYear = 2014 // 10 years before 2024
		yearRange = 10
	} else {
		baseYear = 2029 // 5 years after 2024
		yearRange = 5
	}

	year := baseYear + int(hash[0])%yearRange
	month := int(hash[1])%12 + 1
	day := int(hash[2])%28 + 1 // Use 28 to avoid month-specific day issues

	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// GenerateDeterministicDateOfBirth generates a deterministic date of birth in YYYY-MM-DD format
// Creates realistic ages between 18 and 85 years old (as of 2024)
func GenerateDeterministicDateOfBirth(id, realDOB string) string {
	if realDOB == "" {
		return ""
	}
	hash := hashField(id, "dob", realDOB)

	// Generate ages between 18 and 85 (birth years 1939-2006)
	baseYear := 1939
	yearRange := 68 // 1939 to 2006

	year := baseYear + int(hash[0])%yearRange
	month := int(hash[1])%12 + 1
	day := int(hash[2])%28 + 1 // Use 28 to avoid month-specific day issues

	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// GenerateDeterministicGender generates a deterministic gender
// Returns "Male" or "Female" based on hash
func GenerateDeterministicGender(id, realGender string) string {
	if realGender == "" {
		return ""
	}
	hash := hashField(id, "gender", realGender)

	// Use first byte to determine gender (based on hash of id)
	if int(hash[0])%2 == 0 {
		return "Male"
	}
	return "Female"
}

func GenerateDeterministicSSN(id, realSSN string) string {
	if realSSN == "" {
		return ""
	}
	hash := hashField(id, "ssn", realSSN)

	// Area number (001-899, excluding 000, 666, 900-999)
	areaNum := int(hash[0])%898 + 1 // 1-898

	// Group number (01-99)
	groupNum := int(hash[1])%99 + 1

	// Serial number (0001-9999)
	serialNum := (int(hash[2])<<8 | int(hash[3])) % 9999
	if serialNum == 0 {
		serialNum = 1
	}

	return fmt.Sprintf("%03d-%02d-%04d", areaNum, groupNum, serialNum)
}

// GenerateDeterministicInteger generates a deterministic integer value
// Transforms the input integer in a deterministic way based on ID
func GenerateDeterministicInteger(id string, value int64) int64 {
	if value == 0 {
		return 0
	}

	hash := hashField(id, "integer", fmt.Sprintf("%d", value))

	// Create a deterministic transformation
	// Use hash bytes to create a new integer
	transformed := int64(hash[0])<<56 | int64(hash[1])<<48 | int64(hash[2])<<40 | int64(hash[3])<<32 |
		int64(hash[4])<<24 | int64(hash[5])<<16 | int64(hash[6])<<8 | int64(hash[7])

	// Make it positive and non-zero
	if transformed < 0 {
		transformed = -transformed
	}
	if transformed == 0 {
		transformed = 1
	}

	return transformed
}

// GenerateDeterministicFloat generates a deterministic float value
// Transforms the input float in a deterministic way based on ID
func GenerateDeterministicFloat(id string, value float64) float64 {
	if value == 0 {
		return 0
	}

	hash := hashField(id, "float", fmt.Sprintf("%v", value))

	// Create a deterministic transformation using hash bytes
	intPart := int64(hash[0])<<24 | int64(hash[1])<<16 | int64(hash[2])<<8 | int64(hash[3])
	fracPart := int64(hash[4])<<24 | int64(hash[5])<<16 | int64(hash[6])<<8 | int64(hash[7])

	// Make it positive
	if intPart < 0 {
		intPart = -intPart
	}
	if fracPart < 0 {
		fracPart = -fracPart
	}

	// Normalize fractional part to 0-1 range
	frac := float64(fracPart%1000000) / 1000000.0

	// Add int and fractional parts
	result := float64(intPart%10000000) + frac

	// Ensure it's not zero
	if result == 0 {
		result = 0.1
	}

	return result
}
