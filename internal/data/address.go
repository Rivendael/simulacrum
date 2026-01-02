package data

import (
	"fmt"
	"strings"
)

// ObscureAddress handles address obscuration with support for both combined and separate address fields
func ObscureAddress(id string, real PersonalData) (address, street, city, state, zipCode, county, country string) {
	// Handle address: check if combined address or separate fields are provided
	if real.Address != "" {
		// If combined address provided, use it
		address = GenerateDeterministicAddress(id, real.Address)
	} else if real.Street != "" || real.City != "" || real.State != "" || real.ZipCode != "" {
		// If separate fields provided, only obscure the separate fields (don't generate combined address)
		street = GenerateDeterministicStreet(id, real.Street)
		city = GenerateDeterministicCity(id, real.City)
		state = GenerateDeterministicState(id, real.State)
		zipCode = GenerateDeterministicZipCode(id, real.ZipCode)
	}

	// Always handle county and country if provided
	county = GenerateDeterministicCounty(id, real.County)
	country = GenerateDeterministicCountry(id, real.Country)

	return
}

// GenerateDeterministicAddress generates a deterministic full address
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

// GenerateDeterministicStreet generates a deterministic street address
func GenerateDeterministicStreet(id, realStreet string) string {
	if realStreet == "" {
		return ""
	}
	hash := hashField(id, "street", realStreet)
	number := bytesToInt(hash, 0, 9999)
	street := selectFromList(hash, 2, StreetNames)
	return fmt.Sprintf("%d %s", number, street)
}

// GenerateDeterministicCity generates a deterministic city
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

// GenerateDeterministicState generates a deterministic state
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

// GenerateDeterministicZipCode generates a deterministic zip code
func GenerateDeterministicZipCode(id, realZip string) string {
	if realZip == "" {
		return ""
	}
	hash := hashField(id, "zipcode", realZip)
	zipCode := bytesToInt(hash, 0, 99999)
	return fmt.Sprintf("%05d", zipCode)
}

// GenerateDeterministicCounty generates a deterministic county
func GenerateDeterministicCounty(id, realCounty string) string {
	if realCounty == "" {
		return ""
	}
	hash := hashField(id, "county", realCounty)
	county := selectFromList(hash, 0, CityNames)
	return county
}

// GenerateDeterministicCountry generates a deterministic country
func GenerateDeterministicCountry(id, realCountry string) string {
	if realCountry == "" {
		return ""
	}
	hash := hashField(id, "country", realCountry)
	// Select a region and return its country for consistency
	region := AddressRegions[int(hash[0])%len(AddressRegions)]
	return region.Country
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
