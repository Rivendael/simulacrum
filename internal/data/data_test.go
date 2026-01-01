package data

import (
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
