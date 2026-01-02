package data

import (
	"fmt"
)

// ObscureName handles name obscuration with support for both combined and separate name fields
func ObscureName(id string, real PersonalData) (name, firstName, lastName, middleName string) {
	// Handle name: check if combined name or separate fields are provided
	if real.Name != "" {
		name = GenerateDeterministicName(id, real.Name)
	} else if real.FirstName != "" || real.LastName != "" || real.MiddleName != "" {
		// If separate name fields provided, only obscure the separate fields
		firstName = GenerateDeterministicFirstName(id, real.FirstName)
		lastName = GenerateDeterministicLastName(id, real.LastName)
		middleName = GenerateDeterministicMiddleName(id, real.MiddleName)
	}
	return
}

// GenerateDeterministicName generates a deterministic full name
func GenerateDeterministicName(id, realName string) string {
	if realName == "" {
		return ""
	}
	hash := hashField(id, "name", realName)
	first := selectFromList(hash, 0, FirstNames)
	last := selectFromList(hash, 1, LastNames)
	return fmt.Sprintf("%s %s", first, last)
}

// GenerateDeterministicFirstName generates a deterministic first name
func GenerateDeterministicFirstName(id, realFirstName string) string {
	if realFirstName == "" {
		return ""
	}
	hash := hashField(id, "firstname", realFirstName)
	return selectFromList(hash, 0, FirstNames)
}

// GenerateDeterministicLastName generates a deterministic last name
func GenerateDeterministicLastName(id, realLastName string) string {
	if realLastName == "" {
		return ""
	}
	hash := hashField(id, "lastname", realLastName)
	return selectFromList(hash, 0, LastNames)
}

// GenerateDeterministicMiddleName generates a deterministic middle name
func GenerateDeterministicMiddleName(id, realMiddleName string) string {
	if realMiddleName == "" {
		return ""
	}
	hash := hashField(id, "middlename", realMiddleName)
	return selectFromList(hash, 0, FirstNames)
}
