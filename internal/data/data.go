package data

import (
	"fmt"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type Passport struct {
	Number         string `json:"number"`
	IssueDate      string `json:"issue_date"`
	ExpirationDate string `json:"expiration_date"`
}

type DriverLicense struct {
	Number         string `json:"number"`
	IssueDate      string `json:"issue_date"`
	ExpirationDate string `json:"expiration_date"`
}

type BankAccount struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

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

// buildAddressString combines separate address fields into a single string for obscuration
func buildAddressString(street, city, state, zipCode, county, country string) string {
	parts := []string{}
	if street != "" {
		parts = append(parts, street)
	}
	if city != "" {
		parts = append(parts, city)
	}
	if county != "" {
		parts = append(parts, county)
	}
	if state != "" {
		parts = append(parts, state)
	}
	if zipCode != "" {
		parts = append(parts, zipCode)
	}
	if country != "" {
		parts = append(parts, country)
	}
	return strings.Join(parts, " ")
}

// ObscureData takes the real data and returns a deterministic false identity
func ObscureData(real PersonalData) PersonalData {
	fake := PersonalData{
		ID: real.ID,
	}

	// Handle name: check if combined name or separate fields are provided
	if real.Name != "" {
		fake.Name = GenerateDeterministicName(real.ID, real.Name)
	} else if real.FirstName != "" || real.LastName != "" || real.MiddleName != "" {
		// If separate name fields provided, only obscure the separate fields
		fake.FirstName = GenerateDeterministicFirstName(real.ID, real.FirstName)
		fake.LastName = GenerateDeterministicLastName(real.ID, real.LastName)
		fake.MiddleName = GenerateDeterministicMiddleName(real.ID, real.MiddleName)
	}

	fake.Email = GenerateDeterministicEmail(real.ID, real.Email)

	// Handle address: check if combined address or separate fields are provided
	if real.Address != "" {
		// If combined address provided, use it
		fake.Address = GenerateDeterministicAddress(real.ID, real.Address)
	} else if real.Street != "" || real.City != "" || real.State != "" || real.ZipCode != "" {
		// If separate fields provided, only obscure the separate fields (don't generate combined address)
		fake.Street = GenerateDeterministicStreet(real.ID, real.Street)
		fake.City = GenerateDeterministicCity(real.ID, real.City)
		fake.State = GenerateDeterministicState(real.ID, real.State)
		fake.ZipCode = GenerateDeterministicZipCode(real.ID, real.ZipCode)
	}

	// Always handle county and country if provided
	fake.County = GenerateDeterministicCounty(real.ID, real.County)
	fake.Country = GenerateDeterministicCountry(real.ID, real.Country)

	fake.PhoneNumber = GenerateDeterministicPhone(real.ID, real.PhoneNumber)
	fake.TaxID = GenerateDeterministicTaxID(real.ID, real.TaxID)
	fake.DateOfBirth = GenerateDeterministicDateOfBirth(real.ID, real.DateOfBirth)
	fake.Gender = GenerateDeterministicGender(real.ID, real.Gender)
	fake.SSN = GenerateDeterministicSSN(real.ID, real.SSN)
	if real.Passport != nil {
		fake.Passport = &Passport{
			Number:         GenerateDeterministicPassportNumber(real.ID, real.Passport.Number),
			IssueDate:      GenerateDeterministicDate(real.ID, "passport_issue", real.Passport.IssueDate),
			ExpirationDate: GenerateDeterministicDate(real.ID, "passport_expiry", real.Passport.ExpirationDate),
		}
	}

	// Handle DriverLicense if provided
	if real.DriverLicense != nil {
		fake.DriverLicense = &DriverLicense{
			Number:         GenerateDeterministicDriverLicenseNumber(real.ID, real.DriverLicense.Number),
			IssueDate:      GenerateDeterministicDate(real.ID, "license_issue", real.DriverLicense.IssueDate),
			ExpirationDate: GenerateDeterministicDate(real.ID, "license_expiry", real.DriverLicense.ExpirationDate),
		}
	}

	// Handle BankAccounts if provided
	if len(real.BankAccounts) > 0 {
		fake.BankAccounts = make([]*BankAccount, len(real.BankAccounts))
		for i, account := range real.BankAccounts {
			fake.BankAccounts[i] = &BankAccount{
				Name:   GenerateDeterministicAccountName(real.ID, account.Name, i),
				Amount: GenerateDeterministicAmount(real.ID, account.Amount, i),
			}
		}
	}

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
	for i := 0; i < 8; i++ {
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

func GenerateDeterministicName(id, realName string) string {
	if realName == "" {
		return ""
	}
	hash := hashField(id, "name", realName)
	first := selectFromList(hash, 0, FirstNames)
	last := selectFromList(hash, 1, LastNames)
	return fmt.Sprintf("%s %s", first, last)
}

func GenerateDeterministicFirstName(id, realFirstName string) string {
	if realFirstName == "" {
		return ""
	}
	hash := hashField(id, "firstname", realFirstName)
	return selectFromList(hash, 0, FirstNames)
}

func GenerateDeterministicLastName(id, realLastName string) string {
	if realLastName == "" {
		return ""
	}
	hash := hashField(id, "lastname", realLastName)
	return selectFromList(hash, 0, LastNames)
}

func GenerateDeterministicMiddleName(id, realMiddleName string) string {
	if realMiddleName == "" {
		return ""
	}
	hash := hashField(id, "middlename", realMiddleName)
	return selectFromList(hash, 0, FirstNames)
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

func GenerateDeterministicAddress(id, realAddress string) string {
	if realAddress == "" {
		return ""
	}
	hash := hashField(id, "address", realAddress)
	// Select a region first to ensure consistency
	region := AddressRegions[int(hash[0])%len(AddressRegions)]
	number := bytesToInt(hash, 1, 9999)
	street := selectFromList(hash, 2, StreetNames)
	city := selectFromList(hash, 3, region.Cities)
	state := selectFromList(hash, 4, region.States)
	zipCode := bytesToInt(hash, 5, 99999)
	return fmt.Sprintf("%d %s, %s, %s %05d", number, street, city, state, zipCode)
}

func GenerateDeterministicStreet(id, realStreet string) string {
	if realStreet == "" {
		return ""
	}
	hash := hashField(id, "street", realStreet)
	number := bytesToInt(hash, 0, 9999)
	street := selectFromList(hash, 2, StreetNames)
	return fmt.Sprintf("%d %s", number, street)
}

func GenerateDeterministicCity(id, realCity string) string {
	if realCity == "" {
		return ""
	}
	hash := hashField(id, "city", realCity)
	// Select a region first to ensure consistency
	region := AddressRegions[int(hash[0])%len(AddressRegions)]
	city := selectFromList(hash, 1, region.Cities)
	return city
}

func GenerateDeterministicState(id, realState string) string {
	if realState == "" {
		return ""
	}
	hash := hashField(id, "state", realState)
	// Select a region first to ensure consistency
	region := AddressRegions[int(hash[0])%len(AddressRegions)]
	state := selectFromList(hash, 1, region.States)
	return state
}

func GenerateDeterministicZipCode(id, realZip string) string {
	if realZip == "" {
		return ""
	}
	hash := hashField(id, "zipcode", realZip)
	zipCode := bytesToInt(hash, 0, 99999)
	return fmt.Sprintf("%05d", zipCode)
}

func GenerateDeterministicCounty(id, realCounty string) string {
	if realCounty == "" {
		return ""
	}
	hash := hashField(id, "county", realCounty)
	county := selectFromList(hash, 0, CityNames)
	return county
}

func GenerateDeterministicCountry(id, realCountry string) string {
	if realCountry == "" {
		return ""
	}
	hash := hashField(id, "country", realCountry)
	// Select a region and return its country for consistency
	region := AddressRegions[int(hash[0])%len(AddressRegions)]
	return region.Country
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

// GenerateDeterministicAccountName generates a deterministic bank account name
// Examples: "Checking Account", "Savings Account", "Money Market"
func GenerateDeterministicAccountName(id, realName string, index int) string {
	if realName == "" {
		return ""
	}
	accountTypes := []string{
		"Checking Account",
		"Savings Account",
		"Money Market Account",
		"Deposit Account",
		"Investment Account",
		"Business Account",
		"Interest-Bearing Account",
		"Premium Savings",
		"High-Yield Savings",
		"Retirement Account",
	}

	fieldType := fmt.Sprintf("account_name_%d", index)
	hash := hashField(id, fieldType, realName)
	return selectFromList(hash, 0, accountTypes)
}

// GenerateDeterministicAmount generates a deterministic dollar amount
// Returns amount as a string (e.g., "1234.56")
func GenerateDeterministicAmount(id, realAmount string, index int) string {
	if realAmount == "" {
		return ""
	}

	fieldType := fmt.Sprintf("account_amount_%d", index)
	hash := hashField(id, fieldType, realAmount)

	// Generate realistic amounts: $100 to $999,999.99
	// Use multiple bytes to create a larger number
	cents := (int(hash[0])<<8 | int(hash[1])) % 100
	dollars := (int(hash[2])<<16 | int(hash[3])<<8 | int(hash[4])) % 1000000
	if dollars < 100 {
		dollars += 100
	}

	return fmt.Sprintf("%.2f", float64(dollars)+float64(cents)/100.0)
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
