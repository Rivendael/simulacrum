package data

import (
	"fmt"
	"strings"
	"testing"
)

func TestObscureDataDeterministic(t *testing.T) {
	real := PersonalData{
		ID:          "user123",
		Name:        "John Doe",
		Email:       "john@example.com",
		Address:     "123 Main St",
		PhoneNumber: "555-1234",
	}

	result1 := ObscureData(real)
	result2 := ObscureData(real)

	if result1.Name != result2.Name {
		t.Errorf("Name mismatch: %s != %s", result1.Name, result2.Name)
	}
	if result1.Email != result2.Email {
		t.Errorf("Email mismatch: %s != %s", result1.Email, result2.Email)
	}
	if result1.Address != result2.Address {
		t.Errorf("Address mismatch: %s != %s", result1.Address, result2.Address)
	}
	if result1.PhoneNumber != result2.PhoneNumber {
		t.Errorf("PhoneNumber mismatch: %s != %s", result1.PhoneNumber, result2.PhoneNumber)
	}
}

func TestObscureDataDifferentIDs(t *testing.T) {
	data1 := PersonalData{ID: "user1", Name: "John Doe", Email: "john@example.com"}
	data2 := PersonalData{ID: "user2", Name: "John Doe", Email: "john@example.com"}

	result1 := ObscureData(data1)
	result2 := ObscureData(data2)

	if result1.Name == result2.Name && result1.Email == result2.Email {
		t.Error("Different IDs should produce different obscured data")
	}
}

func TestObscureDataEmptyValues(t *testing.T) {
	real := PersonalData{ID: "user123"}
	result := ObscureData(real)

	if result.ID != "user123" {
		t.Errorf("ID should be preserved, got: %s", result.ID)
	}
	if result.Name != "" || result.Email != "" || result.Address != "" || result.PhoneNumber != "" || result.TaxID != "" {
		t.Error("Empty values should remain empty")
	}
}

func TestGenerateDeterministicTaxID(t *testing.T) {
	taxID1 := GenerateDeterministicTaxID("user123", "123-45-6789")
	taxID2 := GenerateDeterministicTaxID("user123", "123-45-6789")
	taxID3 := GenerateDeterministicTaxID("user124", "123-45-6789")

	if taxID1 != taxID2 {
		t.Errorf("Same ID and TaxID should produce same result: %s != %s", taxID1, taxID2)
	}

	if taxID1 == taxID3 {
		t.Errorf("Different ID should produce different result: %s == %s", taxID1, taxID3)
	}

	if taxID1 == "" || taxID1 == "123-45-6789" {
		t.Errorf("TaxID should be obscured, got: %s", taxID1)
	}

	// Verify format (should be XXX-XX-XXXX)
	if len(taxID1) != 11 || taxID1[3] != '-' || taxID1[6] != '-' {
		t.Errorf("Generated TaxID should match format XXX-XX-XXXX, got: %s", taxID1)
	}
}
func TestObscureDataPreservesID(t *testing.T) {
	testIDs := []string{"cust123", "user_456", "789abc", "test@example.com"}
	for _, id := range testIDs {
		real := PersonalData{ID: id, Name: "Test User"}
		result := ObscureData(real)
		if result.ID != id {
			t.Errorf("ID not preserved: expected %s, got %s", id, result.ID)
		}
	}
}

func TestGenerateDeterministicName(t *testing.T) {
	name1 := GenerateDeterministicName("user123", "John Doe")
	name2 := GenerateDeterministicName("user123", "John Doe")
	name3 := GenerateDeterministicName("user124", "John Doe")

	if name1 != name2 {
		t.Errorf("Same ID and name should produce same result: %s != %s", name1, name2)
	}
	if name1 == name3 {
		t.Errorf("Different ID should produce different result: %s == %s", name1, name3)
	}
	if name1 == "" || name1 == "John Doe" {
		t.Errorf("Name should be obscured, got: %s", name1)
	}
}

func TestGenerateDeterministicNameFormat(t *testing.T) {
	tests := []string{"", "John", "John Doe", "A", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}
	for _, input := range tests {
		result := GenerateDeterministicName("id1", input)
		if input == "" && result != "" {
			t.Errorf("Empty input should produce empty output, got: %s", result)
		}
		if input != "" && result == "" {
			t.Errorf("Non-empty input should produce non-empty output for: %s", input)
		}
		if input != "" && result != "" && !hasExactlyOneSpace(result) {
			t.Errorf("Name should have exactly one space, got: %s", result)
		}
	}
}

// Tests for GenerateDeterministicFirstName
func TestGenerateDeterministicFirstName(t *testing.T) {
	firstName1 := GenerateDeterministicFirstName("user123", "John")
	firstName2 := GenerateDeterministicFirstName("user123", "John")
	firstName3 := GenerateDeterministicFirstName("user124", "John")

	if firstName1 != firstName2 {
		t.Errorf("Same ID and name should produce same result: %s != %s", firstName1, firstName2)
	}
	if firstName1 == firstName3 {
		t.Errorf("Different ID should produce different result: %s == %s", firstName1, firstName3)
	}
	if firstName1 == "" || firstName1 == "John" {
		t.Errorf("FirstName should be obscured, got: %s", firstName1)
	}
}

func TestGenerateDeterministicFirstNameValidation(t *testing.T) {
	validFirstNames := make(map[string]bool)
	for _, name := range FirstNames {
		validFirstNames[name] = true
	}

	for i := 0; i < 50; i++ {
		id := "first_id_" + string(rune(i))
		firstName := GenerateDeterministicFirstName(id, "test firstname")
		if firstName != "" && !validFirstNames[firstName] {
			t.Errorf("Invalid first name: %s", firstName)
		}
	}
}

func TestGenerateDeterministicFirstNameEmpty(t *testing.T) {
	firstName := GenerateDeterministicFirstName("user123", "")
	if firstName != "" {
		t.Errorf("Empty firstName should remain empty, got: %s", firstName)
	}
}

// Tests for GenerateDeterministicLastName
func TestGenerateDeterministicLastName(t *testing.T) {
	lastName1 := GenerateDeterministicLastName("user123", "Doe")
	lastName2 := GenerateDeterministicLastName("user123", "Doe")
	lastName3 := GenerateDeterministicLastName("user124", "Doe")

	if lastName1 != lastName2 {
		t.Errorf("Same ID and name should produce same result: %s != %s", lastName1, lastName2)
	}
	if lastName1 == lastName3 {
		t.Errorf("Different ID should produce different result: %s == %s", lastName1, lastName3)
	}
	if lastName1 == "" || lastName1 == "Doe" {
		t.Errorf("LastName should be obscured, got: %s", lastName1)
	}
}

func TestGenerateDeterministicLastNameValidation(t *testing.T) {
	validLastNames := make(map[string]bool)
	for _, name := range LastNames {
		validLastNames[name] = true
	}

	for i := 0; i < 50; i++ {
		id := "last_id_" + string(rune(i))
		lastName := GenerateDeterministicLastName(id, "test lastname")
		if lastName != "" && !validLastNames[lastName] {
			t.Errorf("Invalid last name: %s", lastName)
		}
	}
}

func TestGenerateDeterministicLastNameEmpty(t *testing.T) {
	lastName := GenerateDeterministicLastName("user123", "")
	if lastName != "" {
		t.Errorf("Empty lastName should remain empty, got: %s", lastName)
	}
}

// Tests for GenerateDeterministicMiddleName
func TestGenerateDeterministicMiddleName(t *testing.T) {
	middleName1 := GenerateDeterministicMiddleName("user123", "John")
	middleName2 := GenerateDeterministicMiddleName("user123", "John")
	middleName3 := GenerateDeterministicMiddleName("user124", "John")

	if middleName1 != middleName2 {
		t.Errorf("Same ID and name should produce same result: %s != %s", middleName1, middleName2)
	}
	if middleName1 == middleName3 {
		t.Errorf("Different ID should produce different result: %s == %s", middleName1, middleName3)
	}
	if middleName1 == "" || middleName1 == "John" {
		t.Errorf("MiddleName should be obscured, got: %s", middleName1)
	}
}

func TestGenerateDeterministicMiddleNameValidation(t *testing.T) {
	validFirstNames := make(map[string]bool)
	for _, name := range FirstNames {
		validFirstNames[name] = true
	}

	for i := 0; i < 50; i++ {
		id := "middle_id_" + string(rune(i))
		middleName := GenerateDeterministicMiddleName(id, "test middlename")
		if middleName != "" && !validFirstNames[middleName] {
			t.Errorf("Invalid middle name: %s", middleName)
		}
	}
}

func TestGenerateDeterministicMiddleNameEmpty(t *testing.T) {
	middleName := GenerateDeterministicMiddleName("user123", "")
	if middleName != "" {
		t.Errorf("Empty middleName should remain empty, got: %s", middleName)
	}
}

// Tests for ObscureData with separate name fields
func TestObscureDataSeparateNameFields(t *testing.T) {
	real := PersonalData{
		ID:         "user456",
		FirstName:  "Jane",
		LastName:   "Smith",
		MiddleName: "Marie",
	}

	result := ObscureData(real)

	if result.ID != "user456" {
		t.Errorf("ID should be preserved, got: %s", result.ID)
	}
	if result.FirstName == "" || result.FirstName == "Jane" {
		t.Errorf("FirstName should be obscured, got: %s", result.FirstName)
	}
	if result.LastName == "" || result.LastName == "Smith" {
		t.Errorf("LastName should be obscured, got: %s", result.LastName)
	}
	if result.MiddleName == "" || result.MiddleName == "Marie" {
		t.Errorf("MiddleName should be obscured, got: %s", result.MiddleName)
	}
	if result.Name != "" {
		t.Errorf("Combined Name should be empty when separate fields are provided, got: %s", result.Name)
	}
}

func TestObscureDataDeterministicSeparateNames(t *testing.T) {
	real := PersonalData{
		ID:         "user789",
		FirstName:  "John",
		LastName:   "Doe",
		MiddleName: "Robert",
	}

	result1 := ObscureData(real)
	result2 := ObscureData(real)

	if result1.FirstName != result2.FirstName {
		t.Errorf("FirstName should be deterministic: %s != %s", result1.FirstName, result2.FirstName)
	}
	if result1.LastName != result2.LastName {
		t.Errorf("LastName should be deterministic: %s != %s", result1.LastName, result2.LastName)
	}
	if result1.MiddleName != result2.MiddleName {
		t.Errorf("MiddleName should be deterministic: %s != %s", result1.MiddleName, result2.MiddleName)
	}
}

func TestObscureDataMixedNameFormats(t *testing.T) {
	// Test 1: Combined name takes priority
	realCombined := PersonalData{
		ID:        "user001",
		Name:      "John Doe",
		FirstName: "Jane",
		LastName:  "Smith",
	}

	resultCombined := ObscureData(realCombined)
	if resultCombined.Name == "" {
		t.Errorf("Combined Name should be populated when Name field is provided")
	}
	if resultCombined.FirstName != "" || resultCombined.LastName != "" {
		t.Errorf("Separate name fields should be empty when combined Name is provided")
	}

	// Test 2: Separate names when combined name is empty
	realSeparate := PersonalData{
		ID:        "user002",
		FirstName: "Alice",
		LastName:  "Johnson",
	}

	resultSeparate := ObscureData(realSeparate)
	if resultSeparate.FirstName == "" || resultSeparate.LastName == "" {
		t.Errorf("Separate name fields should be populated when Name field is empty")
	}
	if resultSeparate.Name != "" {
		t.Errorf("Combined Name should be empty when only separate fields are provided")
	}
}

func TestGenerateDeterministicEmail(t *testing.T) {
	email1 := GenerateDeterministicEmail("user123", "john@example.com")
	email2 := GenerateDeterministicEmail("user123", "john@example.com")
	email3 := GenerateDeterministicEmail("user123", "jane@example.com")

	if email1 != email2 {
		t.Errorf("Same ID and email should produce same result: %s != %s", email1, email2)
	}
	if email1 == email3 {
		t.Errorf("Different email should produce different result: %s == %s", email1, email3)
	}
	if email1 == "" || email1 == "john@example.com" {
		t.Errorf("Email should be obscured, got: %s", email1)
	}
	if !isValidEmail(email1) {
		t.Errorf("Generated email should be valid format, got: %s", email1)
	}
}

func TestGenerateDeterministicEmailDomains(t *testing.T) {
	validDomains := map[string]bool{
		"example.com": true,
		"test.org":    true,
		"fake.net":    true,
		"mail.com":    true,
	}

	for i := 0; i < 256; i++ {
		id := "id_" + string(rune(i))
		email := GenerateDeterministicEmail(id, "test@example.com")
		if email != "" {
			domain := extractDomain(email)
			if !validDomains[domain] {
				t.Errorf("Invalid domain in email: %s", email)
			}
		}
	}
}

func TestGenerateDeterministicAddress(t *testing.T) {
	addr1 := GenerateDeterministicAddress("user123", "123 Main St")
	addr2 := GenerateDeterministicAddress("user123", "123 Main St")
	addr3 := GenerateDeterministicAddress("user124", "123 Main St")

	if addr1 != addr2 {
		t.Errorf("Same ID and address should produce same result: %s != %s", addr1, addr2)
	}
	if addr1 == addr3 {
		t.Errorf("Different ID should produce different result: %s == %s", addr1, addr3)
	}
	if addr1 == "" || addr1 == "123 Main St" {
		t.Errorf("Address should be obscured, got: %s", addr1)
	}
}

func TestGenerateDeterministicAddressFormat(t *testing.T) {
	// Build maps of valid values from lookup lists
	validStreets := make(map[string]bool)
	for _, street := range StreetNames {
		validStreets[street] = true
	}

	// Build maps of all valid cities and states from all regions
	validCities := make(map[string]bool)
	validStates := make(map[string]bool)
	for _, region := range AddressRegions {
		for _, city := range region.Cities {
			validCities[city] = true
		}
		for _, state := range region.States {
			validStates[state] = true
		}
	}

	for i := 0; i < 100; i++ {
		id := "addr_id_" + string(rune(i))
		addr := GenerateDeterministicAddress(id, "test address")
		if addr != "" {
			street := extractStreetFromAddress(addr)
			city := extractCityFromAddress(addr)
			state := extractStateFromAddress(addr)
			zipCode := extractZipFromAddress(addr)

			if !validStreets[street] {
				t.Errorf("Invalid street in address: %s (street: %s)", addr, street)
			}
			if !validCities[city] {
				t.Errorf("Invalid city in address: %s (city: %s)", addr, city)
			}
			if !validStates[state] {
				t.Errorf("Invalid state in address: %s (state: %s)", addr, state)
			}
			if zipCode < 0 || zipCode > 99999 {
				t.Errorf("Invalid zip in address: %s (zip: %d)", addr, zipCode)
			}
		}
	}
}

func TestGenerateDeterministicPhone(t *testing.T) {
	phone1 := GenerateDeterministicPhone("user123", "555-1234")
	phone2 := GenerateDeterministicPhone("user123", "555-1234")
	phone3 := GenerateDeterministicPhone("user124", "555-1234")

	if phone1 != phone2 {
		t.Errorf("Same ID and phone should produce same result: %s != %s", phone1, phone2)
	}
	if phone1 == phone3 {
		t.Errorf("Different ID should produce different result: %s == %s", phone1, phone3)
	}
	if phone1 == "" || phone1 == "555-1234" {
		t.Errorf("Phone should be obscured, got: %s", phone1)
	}
	if !isValidPhoneFormat(phone1) {
		t.Errorf("Generated phone should match format 555-XXX-XXXX, got: %s", phone1)
	}
}

func TestGenerateDeterministicPhoneExchange(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := "phone_id_" + string(rune(i))
		phone := GenerateDeterministicPhone(id, "test phone")
		if phone != "" {
			exchange := extractExchange(phone)
			if exchange < 100 || exchange > 999 {
				t.Errorf("Exchange should be 100-999, got: %d from %s", exchange, phone)
			}
		}
	}
}

func TestGenerateDeterministicTaxIDRanges(t *testing.T) {
	for i := 0; i < 256; i++ {
		id := "tax_id_" + string(rune(i))
		taxID := GenerateDeterministicTaxID(id, "test")
		if taxID != "" {
			area := extractTaxIDArea(taxID)
			group := extractTaxIDGroup(taxID)
			serial := extractTaxIDSerial(taxID)

			if area < 1 || area > 900 {
				t.Errorf("Area should be 1-900, got: %d from %s", area, taxID)
			}
			if group < 1 || group > 100 {
				t.Errorf("Group should be 1-100, got: %d from %s", group, taxID)
			}
			if serial < 0 || serial > 9999 {
				t.Errorf("Serial should be 0-9999, got: %d from %s", serial, taxID)
			}
		}
	}
}

func TestHashFieldConsistency(t *testing.T) {
	hash1 := hashField("id1", "field", "value1")
	hash2 := hashField("id1", "field", "value1")
	hash3 := hashField("id2", "field", "value1")
	hash4 := hashField("id1", "field", "value2")

	if hash1 != hash2 {
		t.Error("Same inputs should produce same hash")
	}
	if hash1 == hash3 {
		t.Error("Different ID should produce different hash")
	}
	if hash1 == hash4 {
		t.Error("Different value should produce different hash")
	}
}

func TestSelectFromListBounds(t *testing.T) {
	lists := [][]string{
		{"A"},
		{"A", "B"},
		{"A", "B", "C", "D", "E"},
		FirstNames,
		LastNames,
	}

	for _, list := range lists {
		for i := 0; i < 256; i++ {
			hash := hashField("id", "test", "val"+string(rune(i)))
			result := selectFromList(hash, 0, list)
			found := false
			for _, item := range list {
				if result == item {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("selectFromList returned item not in list: %s", result)
			}
		}
	}
}

func TestBytesToIntRanges(t *testing.T) {
	tests := []struct {
		max         int
		expectedMin int
		expectedMax int
	}{
		{1, 1, 1},
		{10, 1, 9},
		{100, 1, 99},
		{1000, 1, 999},
		{10000, 1, 9999},
	}

	for _, test := range tests {
		for i := 0; i < 256; i++ {
			hash := hashField("test", "int", string(rune(i)))
			val := bytesToInt(hash, 0, test.max)
			if val < test.expectedMin || val > test.expectedMax {
				t.Errorf("bytesToInt out of range: got %d, expected %d-%d for max %d", val, test.expectedMin, test.expectedMax, test.max)
			}
		}
	}
}

// Helper functions
func isValidEmail(email string) bool {
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			for j := i; j < len(email); j++ {
				if email[j] == '.' {
					return true
				}
			}
		}
	}
	return false
}

func isValidPhoneFormat(phone string) bool {
	if len(phone) != 12 {
		return false
	}
	if phone[0:3] != "555" || phone[3] != '-' || phone[7] != '-' {
		return false
	}
	return true
}

func hasExactlyOneSpace(s string) bool {
	count := 0
	for _, c := range s {
		if c == ' ' {
			count++
		}
	}
	return count == 1
}

func extractDomain(email string) string {
	for i := len(email) - 1; i >= 0; i-- {
		if email[i] == '@' {
			return email[i+1:]
		}
	}
	return ""
}

func extractStreet(addr string) string {
	// Old format - kept for compatibility
	// Address format is: "1234 Street Name Type, City, State Zip"
	// Extract between first space and first comma
	firstSpace := -1
	firstComma := -1
	for i := 0; i < len(addr); i++ {
		if firstSpace == -1 && addr[i] == ' ' {
			firstSpace = i
		}
		if addr[i] == ',' {
			firstComma = i
			break
		}
	}
	if firstSpace >= 0 && firstComma > firstSpace {
		return strings.TrimSpace(addr[firstSpace:firstComma])
	}
	return ""
}

func extractStreetFromAddress(addr string) string {
	// Address format is: "1234 Street Name Type, City, State Zip"
	// Extract between first space and first comma
	firstSpace := -1
	firstComma := -1
	for i := 0; i < len(addr); i++ {
		if firstSpace == -1 && addr[i] == ' ' {
			firstSpace = i
		}
		if addr[i] == ',' {
			firstComma = i
			break
		}
	}
	if firstSpace >= 0 && firstComma > firstSpace {
		return strings.TrimSpace(addr[firstSpace:firstComma])
	}
	return ""
}

func extractCityFromAddress(addr string) string {
	// Address format is: "1234 Street Name Type, City, State Zip"
	// Extract between first comma and second comma
	firstComma := -1
	secondComma := -1
	for i := 0; i < len(addr); i++ {
		if addr[i] == ',' {
			if firstComma == -1 {
				firstComma = i
			} else {
				secondComma = i
				break
			}
		}
	}
	if firstComma >= 0 && secondComma > firstComma {
		return strings.TrimSpace(addr[firstComma+1 : secondComma])
	}
	return ""
}

func extractStateFromAddress(addr string) string {
	// Address format is: "1234 Street Name Type, City, State Zip"
	// Extract between second comma and last space before zip
	secondComma := -1
	commaCount := 0
	for i := 0; i < len(addr); i++ {
		if addr[i] == ',' {
			commaCount++
			if commaCount == 2 {
				secondComma = i
				break
			}
		}
	}
	if secondComma >= 0 {
		// Find the last space after second comma (between state and zip)
		remainder := addr[secondComma+1:]
		for i := len(remainder) - 1; i >= 0; i-- {
			if remainder[i] == ' ' {
				return strings.TrimSpace(remainder[:i])
			}
		}
	}
	return ""
}

func extractZipFromAddress(addr string) int {
	// Address format is: "1234 Street Name Type, City, State Zip"
	// Extract the last 5 digits (zip code)
	if len(addr) < 5 {
		return 0
	}
	lastPart := addr[len(addr)-5:]
	zipCode := 0
	for i := 0; i < 5; i++ {
		if lastPart[i] >= '0' && lastPart[i] <= '9' {
			zipCode = zipCode*10 + int(lastPart[i]-'0')
		} else {
			return 0
		}
	}
	return zipCode
}

func extractExchange(phone string) int {
	if len(phone) < 10 {
		return 0
	}
	exchange := 0
	for i := 4; i < 7; i++ {
		exchange = exchange*10 + int(phone[i]-'0')
	}
	return exchange
}

func extractTaxIDArea(taxID string) int {
	area := 0
	for i := 0; i < 3; i++ {
		area = area*10 + int(taxID[i]-'0')
	}
	return area
}

func extractTaxIDGroup(taxID string) int {
	group := 0
	for i := 4; i < 6; i++ {
		group = group*10 + int(taxID[i]-'0')
	}
	return group
}

func extractTaxIDSerial(taxID string) int {
	// Tax ID format: XXX-XX-XXXX
	// Serial is positions 7-11 (0-indexed)
	if len(taxID) < 11 {
		return 0
	}
	serial := 0
	for i := 7; i < 11; i++ {
		if taxID[i] >= '0' && taxID[i] <= '9' {
			serial = serial*10 + int(taxID[i]-'0')
		}
	}
	return serial
}

// Tests for GenerateDeterministicStreet
func TestGenerateDeterministicStreet(t *testing.T) {
	street1 := GenerateDeterministicStreet("user123", "123 Main St")
	street2 := GenerateDeterministicStreet("user123", "123 Main St")
	street3 := GenerateDeterministicStreet("user124", "123 Main St")

	if street1 != street2 {
		t.Errorf("Same ID and street should produce same result: %s != %s", street1, street2)
	}
	if street1 == street3 {
		t.Errorf("Different ID should produce different result: %s == %s", street1, street3)
	}
	if street1 == "" || street1 == "123 Main St" {
		t.Errorf("Street should be obscured, got: %s", street1)
	}
}

func TestGenerateDeterministicStreetFormat(t *testing.T) {
	validStreets := make(map[string]bool)
	for _, street := range StreetNames {
		validStreets[street] = true
	}

	for i := 0; i < 50; i++ {
		id := "street_id_" + string(rune(i))
		street := GenerateDeterministicStreet(id, "test street")
		if street != "" {
			parts := strings.Split(street, " ")
			if len(parts) < 2 {
				t.Errorf("Street should have number and name: %s", street)
			}
			streetName := strings.Join(parts[1:], " ")
			if !validStreets[streetName] {
				t.Errorf("Invalid street name: %s in %s", streetName, street)
			}
		}
	}
}

func TestGenerateDeterministicStreetEmpty(t *testing.T) {
	street := GenerateDeterministicStreet("user123", "")
	if street != "" {
		t.Errorf("Empty street should remain empty, got: %s", street)
	}
}

// Tests for GenerateDeterministicCity
func TestGenerateDeterministicCity(t *testing.T) {
	city1 := GenerateDeterministicCity("user123", "New York")
	city2 := GenerateDeterministicCity("user123", "New York")
	city3 := GenerateDeterministicCity("user124", "New York")

	if city1 != city2 {
		t.Errorf("Same ID and city should produce same result: %s != %s", city1, city2)
	}
	if city1 == city3 {
		t.Errorf("Different ID should produce different result: %s == %s", city1, city3)
	}
	if city1 == "" || city1 == "New York" {
		t.Errorf("City should be obscured, got: %s", city1)
	}
}

func TestGenerateDeterministicCityValidation(t *testing.T) {
	// Build map of all valid cities from all regions
	validCities := make(map[string]bool)
	for _, region := range AddressRegions {
		for _, city := range region.Cities {
			validCities[city] = true
		}
	}

	for i := 0; i < 50; i++ {
		id := "city_id_" + string(rune(i))
		city := GenerateDeterministicCity(id, "test city")
		if city != "" && !validCities[city] {
			t.Errorf("Invalid city: %s", city)
		}
	}
}

func TestGenerateDeterministicCityEmpty(t *testing.T) {
	city := GenerateDeterministicCity("user123", "")
	if city != "" {
		t.Errorf("Empty city should remain empty, got: %s", city)
	}
}

// Tests for GenerateDeterministicState
func TestGenerateDeterministicState(t *testing.T) {
	state1 := GenerateDeterministicState("user123", "NY")
	state2 := GenerateDeterministicState("user123", "NY")
	state3 := GenerateDeterministicState("user124", "NY")

	if state1 != state2 {
		t.Errorf("Same ID and state should produce same result: %s != %s", state1, state2)
	}
	if state1 == state3 {
		t.Errorf("Different ID should produce different result: %s == %s", state1, state3)
	}
	if state1 == "" || state1 == "NY" {
		t.Errorf("State should be obscured, got: %s", state1)
	}
}

func TestGenerateDeterministicStateValidation(t *testing.T) {
	// Build map of all valid states from all regions
	validStates := make(map[string]bool)
	for _, region := range AddressRegions {
		for _, state := range region.States {
			validStates[state] = true
		}
	}

	for i := 0; i < 100; i++ {
		id := "state_id_" + string(rune(i))
		state := GenerateDeterministicState(id, "test state")
		if state != "" && !validStates[state] {
			t.Errorf("Invalid state: %s", state)
		}
	}
}

func TestGenerateDeterministicStateEmpty(t *testing.T) {
	state := GenerateDeterministicState("user123", "")
	if state != "" {
		t.Errorf("Empty state should remain empty, got: %s", state)
	}
}

// Tests for GenerateDeterministicZipCode
func TestGenerateDeterministicZipCode(t *testing.T) {
	zip1 := GenerateDeterministicZipCode("user123", "12345")
	zip2 := GenerateDeterministicZipCode("user123", "12345")
	zip3 := GenerateDeterministicZipCode("user124", "12345")

	if zip1 != zip2 {
		t.Errorf("Same ID and zip should produce same result: %s != %s", zip1, zip2)
	}
	if zip1 == zip3 {
		t.Errorf("Different ID should produce different result: %s == %s", zip1, zip3)
	}
	if zip1 == "" || zip1 == "12345" {
		t.Errorf("Zip should be obscured, got: %s", zip1)
	}
}

func TestGenerateDeterministicZipCodeFormat(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "zip_id_" + string(rune(i))
		zip := GenerateDeterministicZipCode(id, "test zip")
		if zip != "" {
			if len(zip) != 5 {
				t.Errorf("Zip should be 5 digits, got: %s (length %d)", zip, len(zip))
			}
			for _, c := range zip {
				if c < '0' || c > '9' {
					t.Errorf("Zip should contain only digits, got: %s", zip)
				}
			}
		}
	}
}

func TestGenerateDeterministicZipCodeRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := "zip_range_" + string(rune(i))
		zip := GenerateDeterministicZipCode(id, "test zip")
		if zip != "" {
			zipVal := 0
			for _, c := range zip {
				zipVal = zipVal*10 + int(c-'0')
			}
			if zipVal < 0 || zipVal > 99999 {
				t.Errorf("Zip should be 0-99999, got: %d from %s", zipVal, zip)
			}
		}
	}
}

func TestGenerateDeterministicZipCodeEmpty(t *testing.T) {
	zip := GenerateDeterministicZipCode("user123", "")
	if zip != "" {
		t.Errorf("Empty zip should remain empty, got: %s", zip)
	}
}

// Tests for ObscureData with separate address fields
func TestObscureDataSeparateAddressFields(t *testing.T) {
	real := PersonalData{
		ID:      "user456",
		Name:    "Jane Smith",
		Street:  "456 Oak Ave",
		City:    "Portland",
		State:   "OR",
		ZipCode: "97201",
	}

	result := ObscureData(real)

	if result.ID != "user456" {
		t.Errorf("ID should be preserved, got: %s", result.ID)
	}
	if result.Street == "" || result.Street == "456 Oak Ave" {
		t.Errorf("Street should be obscured, got: %s", result.Street)
	}
	if result.City == "" || result.City == "Portland" {
		t.Errorf("City should be obscured, got: %s", result.City)
	}
	if result.State == "" || result.State == "OR" {
		t.Errorf("State should be obscured, got: %s", result.State)
	}
	if result.ZipCode == "" || result.ZipCode == "97201" {
		t.Errorf("ZipCode should be obscured, got: %s", result.ZipCode)
	}
}

func TestObscureDataDeterministicSeparateFields(t *testing.T) {
	real := PersonalData{
		ID:      "user789",
		Street:  "789 Pine Ln",
		City:    "Seattle",
		State:   "WA",
		ZipCode: "98101",
	}

	result1 := ObscureData(real)
	result2 := ObscureData(real)

	if result1.Street != result2.Street {
		t.Errorf("Street should be deterministic: %s != %s", result1.Street, result2.Street)
	}
	if result1.City != result2.City {
		t.Errorf("City should be deterministic: %s != %s", result1.City, result2.City)
	}
	if result1.State != result2.State {
		t.Errorf("State should be deterministic: %s != %s", result1.State, result2.State)
	}
	if result1.ZipCode != result2.ZipCode {
		t.Errorf("ZipCode should be deterministic: %s != %s", result1.ZipCode, result2.ZipCode)
	}
}

func TestObscureDataMixedAddressFormats(t *testing.T) {
	// Test with combined address only
	real1 := PersonalData{
		ID:      "user1",
		Address: "123 Main St, Boston, MA 02101",
	}
	result1 := ObscureData(real1)
	if result1.Address == "" {
		t.Error("Combined address should be obscured")
	}

	// Test with separate fields only
	real2 := PersonalData{
		ID:      "user2",
		Street:  "456 Oak Ave",
		City:    "Boston",
		State:   "MA",
		ZipCode: "02102",
	}
	result2 := ObscureData(real2)
	if result2.Street == "" || result2.City == "" {
		t.Error("Separate address fields should be obscured")
	}
}

func TestObscureDataPreservesIDWithAddressComponents(t *testing.T) {
	tests := []PersonalData{
		{ID: "id_1", Street: "Street1", City: "City1"},
		{ID: "id_123", City: "City2", State: "ST"},
		{ID: "id_abc_def", Street: "Street3", ZipCode: "12345"},
	}

	for _, test := range tests {
		result := ObscureData(test)
		if result.ID != test.ID {
			t.Errorf("ID should be preserved: expected %s, got %s", test.ID, result.ID)
		}
	}
}

func TestGenerateDeterministicAddressComponentsDifferentIDs(t *testing.T) {
	street1 := GenerateDeterministicStreet("id1", "Main St")
	street2 := GenerateDeterministicStreet("id2", "Main St")
	city1 := GenerateDeterministicCity("id1", "Boston")
	city2 := GenerateDeterministicCity("id2", "Boston")
	state1 := GenerateDeterministicState("id1", "MA")
	state2 := GenerateDeterministicState("id2", "MA")
	zip1 := GenerateDeterministicZipCode("id1", "02101")
	zip2 := GenerateDeterministicZipCode("id2", "02101")

	if street1 == street2 {
		t.Error("Different IDs should produce different streets")
	}
	if city1 == city2 {
		t.Error("Different IDs should produce different cities")
	}
	if state1 == state2 {
		t.Error("Different IDs should produce different states")
	}
	if zip1 == zip2 {
		t.Error("Different IDs should produce different zips")
	}
}

// Tests for GenerateDeterministicCounty
func TestGenerateDeterministicCounty(t *testing.T) {
	county1 := GenerateDeterministicCounty("user123", "test county")
	county2 := GenerateDeterministicCounty("user123", "test county")
	county3 := GenerateDeterministicCounty("user124", "test county")

	if county1 != county2 {
		t.Errorf("Same ID should produce same county: %s != %s", county1, county2)
	}
	if county1 == county3 {
		t.Errorf("Different ID should produce different county: %s == %s", county1, county3)
	}
	if county1 == "" || county1 == "test county" {
		t.Errorf("County should be obscured, got: %s", county1)
	}
}

func TestGenerateDeterministicCountyEmpty(t *testing.T) {
	county := GenerateDeterministicCounty("user123", "")
	if county != "" {
		t.Errorf("Empty county should remain empty, got: %s", county)
	}
}

// Tests for GenerateDeterministicCountry
func TestGenerateDeterministicCountry(t *testing.T) {
	country1 := GenerateDeterministicCountry("user123", "test country")
	country2 := GenerateDeterministicCountry("user123", "test country")
	country3 := GenerateDeterministicCountry("user124", "test country")

	if country1 != country2 {
		t.Errorf("Same ID should produce same country: %s != %s", country1, country2)
	}
	if country1 == country3 {
		t.Errorf("Different ID should produce different country: %s == %s", country1, country3)
	}
	if country1 == "" || country1 == "test country" {
		t.Errorf("Country should be obscured, got: %s", country1)
	}
}

func TestGenerateDeterministicCountryEmpty(t *testing.T) {
	country := GenerateDeterministicCountry("user123", "")
	if country != "" {
		t.Errorf("Empty country should remain empty, got: %s", country)
	}
}

func TestGenerateDeterministicCountryValidation(t *testing.T) {
	// Build map of valid countries from all regions
	validCountries := make(map[string]bool)
	for _, region := range AddressRegions {
		validCountries[region.Country] = true
	}

	for i := 0; i < 50; i++ {
		id := "country_id_" + string(rune(i))
		country := GenerateDeterministicCountry(id, "test country")
		if country != "" && !validCountries[country] {
			t.Errorf("Invalid country: %s", country)
		}
	}
}

// Tests for region-based address consistency
func TestAddressComponentsFromSameRegion(t *testing.T) {
	// When using the full address generator, all components should come from the same region
	id := "region_consistency_test"
	address := GenerateDeterministicAddress(id, "123 Main St, Boston, MA 02101")

	if address == "" {
		t.Skip("Skipping test with empty address")
	}

	// Parse the address to get components
	parts := strings.Split(address, ", ")
	if len(parts) < 3 {
		t.Fatalf("Invalid address format: %s", address)
	}

	city := parts[1]
	stateZip := parts[2]
	stateZipParts := strings.Fields(stateZip)
	if len(stateZipParts) < 1 {
		t.Fatalf("Invalid state/zip format: %s", stateZip)
	}
	state := stateZipParts[0]

	// Find which region these components belong to (they should all be in the same region)
	found := false
	for _, region := range AddressRegions {
		cityFound := false
		stateFound := false
		for _, c := range region.Cities {
			if c == city {
				cityFound = true
				break
			}
		}
		for _, s := range region.States {
			if s == state {
				stateFound = true
				break
			}
		}
		if cityFound && stateFound {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("City %s and State %s do not belong to the same region in address: %s", city, state, address)
	}
}

// Tests for individual city generation from correct region
func TestGenerateDeterministicCityFromRegion(t *testing.T) {
	id := "city_region_test"
	city := GenerateDeterministicCity(id, "Boston")

	if city == "" {
		t.Skip("Skipping test with empty city")
	}

	// Find which region this city belongs to
	hash := hashField(id, "city", "Boston")
	regionIndex := int(hash[0]) % len(AddressRegions)
	expectedRegion := AddressRegions[regionIndex]

	// Verify city belongs to this region
	cityFound := false
	for _, c := range expectedRegion.Cities {
		if c == city {
			cityFound = true
			break
		}
	}
	if !cityFound {
		t.Errorf("City %s not found in its expected region %s", city, expectedRegion.Country)
	}
}

// Tests for individual state generation from correct region
func TestGenerateDeterministicStateFromRegion(t *testing.T) {
	id := "state_region_test"
	state := GenerateDeterministicState(id, "MA")

	if state == "" {
		t.Skip("Skipping test with empty state")
	}

	// Find which region this state belongs to
	hash := hashField(id, "state", "MA")
	regionIndex := int(hash[0]) % len(AddressRegions)
	expectedRegion := AddressRegions[regionIndex]

	// Verify state belongs to this region
	stateFound := false
	for _, s := range expectedRegion.States {
		if s == state {
			stateFound = true
			break
		}
	}
	if !stateFound {
		t.Errorf("State %s not found in its expected region %s", state, expectedRegion.Country)
	}
}

func TestAddressRegionDistribution(t *testing.T) {
	// Test that different IDs can produce addresses from different regions
	regionCounts := make(map[string]int)

	for i := 0; i < 100; i++ {
		id := "region_test_" + string(rune(i%255))
		country := GenerateDeterministicCountry(id, "test")
		if country != "" {
			regionCounts[country]++
		}
	}

	// We should get results from multiple regions (not just one)
	if len(regionCounts) < 2 {
		t.Errorf("Expected addresses from multiple regions, got only %d region(s)", len(regionCounts))
	}
}

func TestAddressComponentsDeterministic(t *testing.T) {
	id := "determinism_test"

	// Each component should be deterministic based on the ID
	results1 := struct {
		city    string
		state   string
		country string
		county  string
	}{
		city:    GenerateDeterministicCity(id, "city"),
		state:   GenerateDeterministicState(id, "state"),
		country: GenerateDeterministicCountry(id, "country"),
		county:  GenerateDeterministicCounty(id, "county"),
	}

	results2 := struct {
		city    string
		state   string
		country string
		county  string
	}{
		city:    GenerateDeterministicCity(id, "city"),
		state:   GenerateDeterministicState(id, "state"),
		country: GenerateDeterministicCountry(id, "country"),
		county:  GenerateDeterministicCounty(id, "county"),
	}

	if results1.city != results2.city {
		t.Errorf("City not deterministic: %s != %s", results1.city, results2.city)
	}
	if results1.state != results2.state {
		t.Errorf("State not deterministic: %s != %s", results1.state, results2.state)
	}
	if results1.country != results2.country {
		t.Errorf("Country not deterministic: %s != %s", results1.country, results2.country)
	}
	if results1.county != results2.county {
		t.Errorf("County not deterministic: %s != %s", results1.county, results2.county)
	}
}

// Tests for GenerateDeterministicPassportNumber
func TestGenerateDeterministicPassportNumber(t *testing.T) {
	passport1 := GenerateDeterministicPassportNumber("user123", "AB1234567")
	passport2 := GenerateDeterministicPassportNumber("user123", "AB1234567")
	passport3 := GenerateDeterministicPassportNumber("user124", "AB1234567")

	if passport1 != passport2 {
		t.Errorf("Same ID and passport should produce same result: %s != %s", passport1, passport2)
	}
	if passport1 == passport3 {
		t.Errorf("Different ID should produce different result: %s == %s", passport1, passport3)
	}
	if passport1 == "" || passport1 == "AB1234567" {
		t.Errorf("Passport should be obscured, got: %s", passport1)
	}
}

func TestGenerateDeterministicPassportNumberFormat(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "passport_id_" + string(rune(i))
		passport := GenerateDeterministicPassportNumber(id, "test passport")
		if passport != "" {
			if len(passport) != 9 {
				t.Errorf("Passport should be 9 chars, got: %s (length %d)", passport, len(passport))
			}
			// First two should be letters
			if passport[0] < 'A' || passport[0] > 'Z' {
				t.Errorf("First char should be A-Z, got: %c", passport[0])
			}
			if passport[1] < 'A' || passport[1] > 'Z' {
				t.Errorf("Second char should be A-Z, got: %c", passport[1])
			}
			// Last 7 should be digits
			for j := 2; j < 9; j++ {
				if passport[j] < '0' || passport[j] > '9' {
					t.Errorf("Char %d should be digit, got: %c in %s", j, passport[j], passport)
				}
			}
		}
	}
}

func TestGenerateDeterministicPassportNumberEmpty(t *testing.T) {
	passport := GenerateDeterministicPassportNumber("user123", "")
	if passport != "" {
		t.Errorf("Empty passport should remain empty, got: %s", passport)
	}
}

// Tests for GenerateDeterministicDriverLicenseNumber
func TestGenerateDeterministicDriverLicenseNumber(t *testing.T) {
	license1 := GenerateDeterministicDriverLicenseNumber("user123", "ABC123456")
	license2 := GenerateDeterministicDriverLicenseNumber("user123", "ABC123456")
	license3 := GenerateDeterministicDriverLicenseNumber("user124", "ABC123456")

	if license1 != license2 {
		t.Errorf("Same ID and license should produce same result: %s != %s", license1, license2)
	}
	if license1 == license3 {
		t.Errorf("Different ID should produce different result: %s == %s", license1, license3)
	}
	if license1 == "" || license1 == "ABC123456" {
		t.Errorf("Driver License should be obscured, got: %s", license1)
	}
}

func TestGenerateDeterministicDriverLicenseNumberFormat(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "license_id_" + string(rune(i))
		license := GenerateDeterministicDriverLicenseNumber(id, "test license")
		if license != "" {
			if len(license) != 9 {
				t.Errorf("License should be 9 chars, got: %s (length %d)", license, len(license))
			}
			// First three should be letters
			for j := 0; j < 3; j++ {
				if license[j] < 'A' || license[j] > 'Z' {
					t.Errorf("Char %d should be A-Z, got: %c", j, license[j])
				}
			}
			// Last 6 should be digits
			for j := 3; j < 9; j++ {
				if license[j] < '0' || license[j] > '9' {
					t.Errorf("Char %d should be digit, got: %c in %s", j, license[j], license)
				}
			}
		}
	}
}

func TestGenerateDeterministicDriverLicenseNumberEmpty(t *testing.T) {
	license := GenerateDeterministicDriverLicenseNumber("user123", "")
	if license != "" {
		t.Errorf("Empty license should remain empty, got: %s", license)
	}
}

// Tests for GenerateDeterministicDate
func TestGenerateDeterministicDate(t *testing.T) {
	date1 := GenerateDeterministicDate("user123", "passport_issue", "2020-01-15")
	date2 := GenerateDeterministicDate("user123", "passport_issue", "2020-01-15")
	date3 := GenerateDeterministicDate("user124", "passport_issue", "2020-01-15")

	if date1 != date2 {
		t.Errorf("Same ID and date should produce same result: %s != %s", date1, date2)
	}
	if date1 == date3 {
		t.Errorf("Different ID should produce different result: %s == %s", date1, date3)
	}
	if date1 == "" || date1 == "2020-01-15" {
		t.Errorf("Date should be obscured, got: %s", date1)
	}
}

func TestGenerateDeterministicDateFormat(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "date_id_" + string(rune(i))
		date := GenerateDeterministicDate(id, "passport_issue", "test date")
		if date != "" {
			if len(date) != 10 {
				t.Errorf("Date should be 10 chars (YYYY-MM-DD), got: %s (length %d)", date, len(date))
			}
			// Check format: YYYY-MM-DD
			if date[4] != '-' || date[7] != '-' {
				t.Errorf("Date should have dashes at positions 4 and 7, got: %s", date)
			}
			// Validate year range for issue dates (2014-2024)
			if date[:4] < "2014" || date[:4] > "2024" {
				t.Errorf("Issue date year should be 2014-2024, got: %s", date[:4])
			}
		}
	}
}

func TestGenerateDeterministicDateExpirationRange(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "expiry_id_" + string(rune(i))
		date := GenerateDeterministicDate(id, "passport_expiry", "test date")
		if date != "" {
			// Validate year range for expiration dates (2029-2034)
			if date[:4] < "2029" || date[:4] > "2034" {
				t.Errorf("Expiration date year should be 2029-2034, got: %s", date[:4])
			}
		}
	}
}

func TestGenerateDeterministicDateEmpty(t *testing.T) {
	date := GenerateDeterministicDate("user123", "passport_issue", "")
	if date != "" {
		t.Errorf("Empty date should remain empty, got: %s", date)
	}
}

// Tests for ObscureData with Passport
func TestObscureDataWithPassport(t *testing.T) {
	real := PersonalData{
		ID: "user456",
		Passport: &Passport{
			Number:         "AB1234567",
			IssueDate:      "2020-01-15",
			ExpirationDate: "2030-01-15",
		},
	}

	result := ObscureData(real)

	if result.ID != "user456" {
		t.Errorf("ID should be preserved, got: %s", result.ID)
	}
	if result.Passport == nil {
		t.Errorf("Passport should not be nil")
	} else {
		if result.Passport.Number == "" || result.Passport.Number == "AB1234567" {
			t.Errorf("Passport number should be obscured, got: %s", result.Passport.Number)
		}
		if result.Passport.IssueDate == "" || result.Passport.IssueDate == "2020-01-15" {
			t.Errorf("Passport issue date should be obscured, got: %s", result.Passport.IssueDate)
		}
		if result.Passport.ExpirationDate == "" || result.Passport.ExpirationDate == "2030-01-15" {
			t.Errorf("Passport expiration date should be obscured, got: %s", result.Passport.ExpirationDate)
		}
	}
}

func TestObscureDataWithDriverLicense(t *testing.T) {
	real := PersonalData{
		ID: "user789",
		DriverLicense: &DriverLicense{
			Number:         "DL123456",
			IssueDate:      "2019-06-10",
			ExpirationDate: "2029-06-10",
		},
	}

	result := ObscureData(real)

	if result.ID != "user789" {
		t.Errorf("ID should be preserved, got: %s", result.ID)
	}
	if result.DriverLicense == nil {
		t.Errorf("DriverLicense should not be nil")
	} else {
		if result.DriverLicense.Number == "" || result.DriverLicense.Number == "DL123456" {
			t.Errorf("Driver license number should be obscured, got: %s", result.DriverLicense.Number)
		}
		if result.DriverLicense.IssueDate == "" || result.DriverLicense.IssueDate == "2019-06-10" {
			t.Errorf("Driver license issue date should be obscured, got: %s", result.DriverLicense.IssueDate)
		}
		if result.DriverLicense.ExpirationDate == "" || result.DriverLicense.ExpirationDate == "2029-06-10" {
			t.Errorf("Driver license expiration date should be obscured, got: %s", result.DriverLicense.ExpirationDate)
		}
	}
}

func TestObscureDataDeterministicPassportAndLicense(t *testing.T) {
	real := PersonalData{
		ID: "user999",
		Passport: &Passport{
			Number:         "XY9876543",
			IssueDate:      "2018-05-20",
			ExpirationDate: "2028-05-20",
		},
		DriverLicense: &DriverLicense{
			Number:         "LIC999999",
			IssueDate:      "2018-05-20",
			ExpirationDate: "2028-05-20",
		},
	}

	result1 := ObscureData(real)
	result2 := ObscureData(real)

	if result1.Passport.Number != result2.Passport.Number {
		t.Errorf("Passport number should be deterministic: %s != %s", result1.Passport.Number, result2.Passport.Number)
	}
	if result1.Passport.IssueDate != result2.Passport.IssueDate {
		t.Errorf("Passport issue date should be deterministic: %s != %s", result1.Passport.IssueDate, result2.Passport.IssueDate)
	}
	if result1.DriverLicense.Number != result2.DriverLicense.Number {
		t.Errorf("Driver license number should be deterministic: %s != %s", result1.DriverLicense.Number, result2.DriverLicense.Number)
	}
	if result1.DriverLicense.IssueDate != result2.DriverLicense.IssueDate {
		t.Errorf("Driver license issue date should be deterministic: %s != %s", result1.DriverLicense.IssueDate, result2.DriverLicense.IssueDate)
	}
}

func TestObscureDataWithoutPassportAndLicense(t *testing.T) {
	real := PersonalData{
		ID:    "user000",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	result := ObscureData(real)

	if result.Passport != nil {
		t.Errorf("Passport should be nil when not provided, got: %v", result.Passport)
	}
	if result.DriverLicense != nil {
		t.Errorf("DriverLicense should be nil when not provided, got: %v", result.DriverLicense)
	}
}

// Tests for GenerateDeterministicDateOfBirth
func TestGenerateDeterministicDateOfBirth(t *testing.T) {
	dob1 := GenerateDeterministicDateOfBirth("user123", "1990-05-15")
	dob2 := GenerateDeterministicDateOfBirth("user123", "1990-05-15")
	dob3 := GenerateDeterministicDateOfBirth("user124", "1990-05-15")

	if dob1 != dob2 {
		t.Errorf("Same ID and DOB should produce same result: %s != %s", dob1, dob2)
	}
	if dob1 == dob3 {
		t.Errorf("Different ID should produce different result: %s == %s", dob1, dob3)
	}
	if dob1 == "" || dob1 == "1990-05-15" {
		t.Errorf("DOB should be obscured, got: %s", dob1)
	}
}

func TestGenerateDeterministicDateOfBirthFormat(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "dob_id_" + string(rune(i))
		dob := GenerateDeterministicDateOfBirth(id, "test dob")
		if dob != "" {
			if len(dob) != 10 {
				t.Errorf("DOB should be 10 chars (YYYY-MM-DD), got: %s (length %d)", dob, len(dob))
			}
			// Check format: YYYY-MM-DD
			if dob[4] != '-' || dob[7] != '-' {
				t.Errorf("DOB should have dashes at positions 4 and 7, got: %s", dob)
			}
			// Validate year range (1939-2006 for ages 18-85)
			if dob[:4] < "1939" || dob[:4] > "2006" {
				t.Errorf("DOB year should be 1939-2006, got: %s", dob[:4])
			}
		}
	}
}

func TestGenerateDeterministicDateOfBirthEmpty(t *testing.T) {
	dob := GenerateDeterministicDateOfBirth("user123", "")
	if dob != "" {
		t.Errorf("Empty DOB should remain empty, got: %s", dob)
	}
}

// Tests for GenerateDeterministicGender
func TestGenerateDeterministicGender(t *testing.T) {
	gender1 := GenerateDeterministicGender("user123", "M")
	gender2 := GenerateDeterministicGender("user123", "M")

	if gender1 != gender2 {
		t.Errorf("Same ID should produce same result: %s != %s", gender1, gender2)
	}
	if gender1 == "" || (gender1 != "Male" && gender1 != "Female") {
		t.Errorf("Gender should be Male or Female, got: %s", gender1)
	}
}

func TestGenerateDeterministicGenderValid(t *testing.T) {
	maleCount := 0
	femaleCount := 0

	for i := 0; i < 100; i++ {
		id := "gender_id_" + string(rune(i))
		gender := GenerateDeterministicGender(id, "test gender")
		if gender == "Male" {
			maleCount++
		} else if gender == "Female" {
			femaleCount++
		} else {
			t.Errorf("Invalid gender: %s", gender)
		}
	}

	// Should have a reasonable distribution
	if maleCount == 0 || femaleCount == 0 {
		t.Errorf("Gender distribution skewed: Male=%d, Female=%d", maleCount, femaleCount)
	}
}

func TestGenerateDeterministicGenderEmpty(t *testing.T) {
	gender := GenerateDeterministicGender("user123", "")
	if gender != "" {
		t.Errorf("Empty gender should remain empty, got: %s", gender)
	}
}

// Tests for GenerateDeterministicSSN
func TestGenerateDeterministicSSN(t *testing.T) {
	ssn1 := GenerateDeterministicSSN("user123", "123-45-6789")
	ssn2 := GenerateDeterministicSSN("user123", "123-45-6789")
	ssn3 := GenerateDeterministicSSN("user124", "123-45-6789")

	if ssn1 != ssn2 {
		t.Errorf("Same ID and SSN should produce same result: %s != %s", ssn1, ssn2)
	}
	if ssn1 == ssn3 {
		t.Errorf("Different ID should produce different result: %s == %s", ssn1, ssn3)
	}
	if ssn1 == "" || ssn1 == "123-45-6789" {
		t.Errorf("SSN should be obscured, got: %s", ssn1)
	}
}

func TestGenerateDeterministicSSNFormat(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "ssn_id_" + string(rune(i))
		ssn := GenerateDeterministicSSN(id, "test ssn")
		if ssn != "" {
			if len(ssn) != 11 {
				t.Errorf("SSN should be 11 chars (XXX-XX-XXXX), got: %s (length %d)", ssn, len(ssn))
			}
			// Check format: XXX-XX-XXXX
			if ssn[3] != '-' || ssn[6] != '-' {
				t.Errorf("SSN should have dashes at positions 3 and 6, got: %s", ssn)
			}
			// Validate all other characters are digits
			for j := 0; j < len(ssn); j++ {
				if j == 3 || j == 6 {
					continue
				}
				if ssn[j] < '0' || ssn[j] > '9' {
					t.Errorf("SSN should contain only digits and hyphens, got: %c in %s", ssn[j], ssn)
				}
			}
			// Area number should be 1-898 (001-898)
			areaStr := ssn[:3]
			if areaStr < "001" || areaStr > "898" {
				t.Errorf("SSN area should be 001-898, got: %s in %s", areaStr, ssn)
			}
		}
	}
}

func TestGenerateDeterministicSSNEmpty(t *testing.T) {
	ssn := GenerateDeterministicSSN("user123", "")
	if ssn != "" {
		t.Errorf("Empty SSN should remain empty, got: %s", ssn)
	}
}

// Tests for ObscureData with new demographic fields
func TestObscureDataWithDemographics(t *testing.T) {
	real := PersonalData{
		ID:          "user_demo_001",
		Name:        "Patricia Lee",
		Email:       "patricia.lee@company.com",
		DateOfBirth: "1985-07-22",
		Gender:      "Female",
		SSN:         "123-45-6789",
	}

	result := ObscureData(real)

	if result.ID != "user_demo_001" {
		t.Errorf("ID should be preserved, got: %s", result.ID)
	}
	if result.DateOfBirth == "" || result.DateOfBirth == "1985-07-22" {
		t.Errorf("DateOfBirth should be obscured, got: %s", result.DateOfBirth)
	}
	// Gender might be same or different - both are valid
	if result.Gender == "" || (result.Gender != "Male" && result.Gender != "Female") {
		t.Errorf("Gender should be Male or Female, got: %s", result.Gender)
	}
	if result.SSN == "" || result.SSN == "123-45-6789" {
		t.Errorf("SSN should be obscured, got: %s", result.SSN)
	}
}

func TestObscureDataDeterministicDemographics(t *testing.T) {
	real := PersonalData{
		ID:          "user_demo_002",
		DateOfBirth: "1992-03-10",
		Gender:      "Male",
		SSN:         "987-65-4321",
	}

	result1 := ObscureData(real)
	result2 := ObscureData(real)

	if result1.DateOfBirth != result2.DateOfBirth {
		t.Errorf("DateOfBirth should be deterministic: %s != %s", result1.DateOfBirth, result2.DateOfBirth)
	}
	if result1.Gender != result2.Gender {
		t.Errorf("Gender should be deterministic: %s != %s", result1.Gender, result2.Gender)
	}
	if result1.SSN != result2.SSN {
		t.Errorf("SSN should be deterministic: %s != %s", result1.SSN, result2.SSN)
	}
}

func TestObscureDataWithoutDemographics(t *testing.T) {
	real := PersonalData{
		ID:    "user_no_demo",
		Name:  "John Smith",
		Email: "john.smith@company.com",
	}

	result := ObscureData(real)

	if result.DateOfBirth != "" {
		t.Errorf("DateOfBirth should be empty when not provided, got: %s", result.DateOfBirth)
	}
	if result.Gender != "" {
		t.Errorf("Gender should be empty when not provided, got: %s", result.Gender)
	}
	if result.SSN != "" {
		t.Errorf("SSN should be empty when not provided, got: %s", result.SSN)
	}
}

// Tests for GenerateDeterministicAccountName
func TestGenerateDeterministicAccountName(t *testing.T) {
	name1 := GenerateDeterministicAccountName("user123", "Checking Account", 0)
	name2 := GenerateDeterministicAccountName("user123", "Checking Account", 0)
	name3 := GenerateDeterministicAccountName("user124", "Checking Account", 0)

	if name1 != name2 {
		t.Errorf("Same ID and index should produce same result: %s != %s", name1, name2)
	}
	if name1 == "" || name1 == "Checking Account" {
		t.Errorf("Account name should be obscured, got: %s", name1)
	}
	if name1 == name3 {
		t.Errorf("Different ID should produce different result: %s == %s", name1, name3)
	}
}

func TestGenerateDeterministicAccountNameValid(t *testing.T) {
	validNames := map[string]bool{
		"Checking Account":         true,
		"Savings Account":          true,
		"Money Market Account":     true,
		"Deposit Account":          true,
		"Investment Account":       true,
		"Business Account":         true,
		"Interest-Bearing Account": true,
		"Premium Savings":          true,
		"High-Yield Savings":       true,
		"Retirement Account":       true,
	}

	for i := 0; i < 20; i++ {
		id := "account_name_" + string(rune(i))
		name := GenerateDeterministicAccountName(id, "test account", 0)
		if name != "" && !validNames[name] {
			t.Errorf("Invalid account name: %s", name)
		}
	}
}

func TestGenerateDeterministicAccountNameEmpty(t *testing.T) {
	name := GenerateDeterministicAccountName("user123", "", 0)
	if name != "" {
		t.Errorf("Empty account name should remain empty, got: %s", name)
	}
}

func TestGenerateDeterministicAccountNameIndex(t *testing.T) {
	// Different indices should produce different results
	name0 := GenerateDeterministicAccountName("user123", "Checking Account", 0)
	name1 := GenerateDeterministicAccountName("user123", "Checking Account", 1)
	name2 := GenerateDeterministicAccountName("user123", "Checking Account", 2)

	// Should be able to get different accounts for different indices
	accountSet := make(map[string]bool)
	accountSet[name0] = true
	accountSet[name1] = true
	accountSet[name2] = true

	if len(accountSet) == 1 {
		t.Errorf("Different indices should produce some different results: %s, %s, %s", name0, name1, name2)
	}
}

// Tests for GenerateDeterministicAmount
func TestGenerateDeterministicAmount(t *testing.T) {
	amount1 := GenerateDeterministicAmount("user123", "1000.00", 0)
	amount2 := GenerateDeterministicAmount("user123", "1000.00", 0)
	amount3 := GenerateDeterministicAmount("user124", "1000.00", 0)

	if amount1 != amount2 {
		t.Errorf("Same ID and index should produce same result: %s != %s", amount1, amount2)
	}
	if amount1 == "" || amount1 == "1000.00" {
		t.Errorf("Amount should be obscured, got: %s", amount1)
	}
	if amount1 == amount3 {
		t.Errorf("Different ID should produce different result: %s == %s", amount1, amount3)
	}
}

func TestGenerateDeterministicAmountFormat(t *testing.T) {
	for i := 0; i < 20; i++ {
		id := "amount_id_" + string(rune(i))
		amount := GenerateDeterministicAmount(id, "5000.00", 0)
		if amount != "" {
			// Check format: should be numeric with decimal
			if !strings.Contains(amount, ".") {
				t.Errorf("Amount should contain decimal point, got: %s", amount)
			}
			// Should have exactly 2 decimal places
			parts := strings.Split(amount, ".")
			if len(parts) != 2 || len(parts[1]) != 2 {
				t.Errorf("Amount should have 2 decimal places, got: %s", amount)
			}
		}
	}
}

func TestGenerateDeterministicAmountRange(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := "amount_range_" + string(rune(i))
		amountStr := GenerateDeterministicAmount(id, "test amount", 0)
		if amountStr != "" {
			var amount float64
			_, err := fmt.Sscanf(amountStr, "%f", &amount)
			if err != nil {
				t.Errorf("Amount should be parseable as float: %s", amountStr)
			}
			if amount < 100 || amount > 1000000 {
				t.Errorf("Amount should be between 100 and 1000000, got: %f", amount)
			}
		}
	}
}

func TestGenerateDeterministicAmountEmpty(t *testing.T) {
	amount := GenerateDeterministicAmount("user123", "", 0)
	if amount != "" {
		t.Errorf("Empty amount should remain empty, got: %s", amount)
	}
}

// Tests for ObscureData with BankAccounts
func TestObscureDataWithBankAccounts(t *testing.T) {
	real := PersonalData{
		ID:   "user_bank_001",
		Name: "David Chen",
		BankAccounts: []*BankAccount{
			{Name: "Checking Account", Amount: "5000.00"},
			{Name: "Savings Account", Amount: "25000.50"},
			{Name: "Money Market Account", Amount: "75000.00"},
		},
	}

	result := ObscureData(real)

	if result.ID != "user_bank_001" {
		t.Errorf("ID should be preserved, got: %s", result.ID)
	}
	if len(result.BankAccounts) != 3 {
		t.Errorf("Should have 3 bank accounts, got: %d", len(result.BankAccounts))
	}

	for i, account := range result.BankAccounts {
		if account.Name == "" {
			t.Errorf("Account %d name should not be empty", i)
		}
		if account.Amount == "" {
			t.Errorf("Account %d amount should not be empty", i)
		}
		if account.Name == real.BankAccounts[i].Name {
			t.Errorf("Account %d name should be obscured", i)
		}
		if account.Amount == real.BankAccounts[i].Amount {
			t.Errorf("Account %d amount should be obscured", i)
		}
	}
}

func TestObscureDataDeterministicBankAccounts(t *testing.T) {
	real := PersonalData{
		ID: "user_bank_002",
		BankAccounts: []*BankAccount{
			{Name: "Checking Account", Amount: "3000.00"},
			{Name: "Savings Account", Amount: "50000.00"},
		},
	}

	result1 := ObscureData(real)
	result2 := ObscureData(real)

	if len(result1.BankAccounts) != len(result2.BankAccounts) {
		t.Errorf("Number of accounts should be deterministic")
	}

	for i := 0; i < len(result1.BankAccounts); i++ {
		if result1.BankAccounts[i].Name != result2.BankAccounts[i].Name {
			t.Errorf("Account %d name should be deterministic: %s != %s", i, result1.BankAccounts[i].Name, result2.BankAccounts[i].Name)
		}
		if result1.BankAccounts[i].Amount != result2.BankAccounts[i].Amount {
			t.Errorf("Account %d amount should be deterministic: %s != %s", i, result1.BankAccounts[i].Amount, result2.BankAccounts[i].Amount)
		}
	}
}

func TestObscureDataWithoutBankAccounts(t *testing.T) {
	real := PersonalData{
		ID:    "user_no_bank",
		Name:  "Lisa Park",
		Email: "lisa.park@company.com",
	}

	result := ObscureData(real)

	if result.BankAccounts != nil && len(result.BankAccounts) > 0 {
		t.Errorf("BankAccounts should be empty when not provided, got: %v", result.BankAccounts)
	}
}

func TestObscureDataEmptyBankAccounts(t *testing.T) {
	real := PersonalData{
		ID:           "user_empty_bank",
		Name:         "Maria Garcia",
		BankAccounts: []*BankAccount{},
	}

	result := ObscureData(real)

	if result.BankAccounts != nil && len(result.BankAccounts) > 0 {
		t.Errorf("BankAccounts should be empty, got: %v", result.BankAccounts)
	}
}

func TestGenerateDeterministicInteger(t *testing.T) {
	id := "user123"
	value := int64(42)
	result := GenerateDeterministicInteger(id, value)

	if result == 0 {
		t.Errorf("Expected non-zero integer")
	}
	if result < 0 {
		t.Errorf("Expected positive integer, got %d", result)
	}
}

func TestGenerateDeterministicIntegerDeterminism(t *testing.T) {
	id := "user123"
	value := int64(42)
	result1 := GenerateDeterministicInteger(id, value)
	result2 := GenerateDeterministicInteger(id, value)

	if result1 != result2 {
		t.Errorf("Expected deterministic results: %d != %d", result1, result2)
	}
}

func TestGenerateDeterministicIntegerZero(t *testing.T) {
	id := "user123"
	result := GenerateDeterministicInteger(id, 0)

	if result != 0 {
		t.Errorf("Expected 0 for zero input, got %d", result)
	}
}

func TestGenerateDeterministicIntegerDifferentIDs(t *testing.T) {
	value := int64(42)
	result1 := GenerateDeterministicInteger("user1", value)
	result2 := GenerateDeterministicInteger("user2", value)

	if result1 == result2 {
		t.Errorf("Expected different results for different IDs")
	}
}

func TestGenerateDeterministicFloat(t *testing.T) {
	id := "user123"
	value := 123.45

	result := GenerateDeterministicFloat(id, value)

	if result == 0 {
		t.Errorf("Expected non-zero float")
	}
	if result < 0 {
		t.Errorf("Expected positive float, got %f", result)
	}
}

func TestGenerateDeterministicFloatDeterminism(t *testing.T) {
	id := "user123"
	value := 123.45
	result1 := GenerateDeterministicFloat(id, value)
	result2 := GenerateDeterministicFloat(id, value)

	if result1 != result2 {
		t.Errorf("Expected deterministic results: %f != %f", result1, result2)
	}
}

func TestGenerateDeterministicFloatZero(t *testing.T) {
	id := "user123"
	result := GenerateDeterministicFloat(id, 0)

	if result != 0 {
		t.Errorf("Expected 0 for zero input, got %f", result)
	}
}

func TestGenerateDeterministicFloatDifferentIDs(t *testing.T) {
	value := 123.45
	result1 := GenerateDeterministicFloat("user1", value)
	result2 := GenerateDeterministicFloat("user2", value)

	if result1 == result2 {
		t.Errorf("Expected different results for different IDs")
	}
}

func TestObscureDataWithIntegerAndFloat(t *testing.T) {
	input := PersonalData{
		ID:           "user123",
		Name:         "John Doe",
		IntegerValue: 12345,
		FloatValue:   678.90,
	}

	result := ObscureData(input)

	if result.ID != input.ID {
		t.Errorf("Expected ID to match")
	}
	if result.IntegerValue == 0 {
		t.Errorf("Expected non-zero integer value")
	}
	if result.FloatValue == 0 {
		t.Errorf("Expected non-zero float value")
	}
}

func TestObscureDataWithoutIntegerAndFloat(t *testing.T) {
	input := PersonalData{
		ID:   "user123",
		Name: "John Doe",
	}

	result := ObscureData(input)

	if result.IntegerValue != 0 {
		t.Errorf("Expected zero integer value for empty input")
	}
	if result.FloatValue != 0 {
		t.Errorf("Expected zero float value for empty input")
	}
}

func TestObscureDataIntegerFloatDeterminism(t *testing.T) {
	input := PersonalData{
		ID:           "user123",
		Name:         "John Doe",
		IntegerValue: 12345,
		FloatValue:   678.90,
	}

	result1 := ObscureData(input)
	result2 := ObscureData(input)

	if result1.IntegerValue != result2.IntegerValue {
		t.Errorf("Expected deterministic integer results")
	}
	if result1.FloatValue != result2.FloatValue {
		t.Errorf("Expected deterministic float results")
	}
}

func TestGenerateDeterministicPhoneWithCountryCode(t *testing.T) {
	phone := "+1-555-987-6543"
	result := GenerateDeterministicPhone("user123", phone)

	if result == phone {
		t.Errorf("Expected phone to be obscured")
	}

	if !strings.HasPrefix(result, "+1-") {
		t.Errorf("Expected country code +1- to be preserved, got %s", result)
	}

	// Should be deterministic
	result2 := GenerateDeterministicPhone("user123", phone)
	if result != result2 {
		t.Errorf("Expected deterministic results: %s vs %s", result, result2)
	}
}

func TestGenerateDeterministicPhoneWithoutCountryCode(t *testing.T) {
	phone := "555-987-6543"
	result := GenerateDeterministicPhone("user123", phone)

	if result == phone {
		t.Errorf("Expected phone to be obscured")
	}

	if strings.HasPrefix(result, "+") {
		t.Errorf("Expected no country code, got %s", result)
	}

	// Should be deterministic
	result2 := GenerateDeterministicPhone("user123", phone)
	if result != result2 {
		t.Errorf("Expected deterministic results: %s vs %s", result, result2)
	}
}

func TestGenerateDeterministicPhoneDifferentCountryCodes(t *testing.T) {
	phone1 := "+1-555-123-4567"
	phone2 := "+44-555-123-4567" // Different country code

	result1 := GenerateDeterministicPhone("user123", phone1)
	result2 := GenerateDeterministicPhone("user123", phone2)

	if result1 == result2 {
		t.Errorf("Expected different results for different country codes")
	}

	if !strings.HasPrefix(result1, "+1-") {
		t.Errorf("Expected result1 to have +1-, got %s", result1)
	}

	if !strings.HasPrefix(result2, "+44-") {
		t.Errorf("Expected result2 to have +44-, got %s", result2)
	}
}

func TestGenerateDeterministicPhoneEmpty(t *testing.T) {
	result := GenerateDeterministicPhone("user123", "")
	if result != "" {
		t.Errorf("Expected empty string for empty input, got %s", result)
	}
}
