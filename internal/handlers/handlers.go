package handlers

import (
	"fmt"
	"io"
	"net/http"

	"simulacrum/internal/data"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
)

// List of fields that should be obscured
var obscurableFields = map[string]bool{
	"id":             true,
	"name":           true,
	"email":          true,
	"phone_number":   true,
	"address":        true,
	"street":         true,
	"city":           true,
	"state":          true,
	"zip_code":       true,
	"county":         true,
	"country":        true,
	"tax_id":         true,
	"first_name":     true,
	"last_name":      true,
	"middle_name":    true,
	"date_of_birth":  true,
	"gender":         true,
	"ssn":            true,
	"passport":       true,
	"driver_license": true,
	"bank_accounts":  true,
	"integer_value":  true,
	"float_value":    true,
}

// HandleObscure accepts arbitrary JSON and obscures any recognized fields
func HandleObscure(c *gin.Context) {
	// Use fastjson for faster unmarshaling
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	var parser fastjson.Parser
	v, err := parser.ParseBytes(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Convert fastjson Value to map[string]interface{}
	req := fastjsonToInterface(v)
	result := obscureGeneric(req)
	c.JSON(http.StatusOK, result)
}

// fastjsonToInterface converts a fastjson.Value to interface{}
func fastjsonToInterface(v *fastjson.Value) interface{} {
	switch v.Type() {
	case fastjson.TypeObject:
		m := make(map[string]interface{})
		obj, _ := v.Object()
		obj.Visit(func(key []byte, v *fastjson.Value) {
			m[string(key)] = fastjsonToInterface(v)
		})
		return m
	case fastjson.TypeArray:
		arr := v.GetArray()
		result := make([]interface{}, len(arr))
		for i, item := range arr {
			result[i] = fastjsonToInterface(item)
		}
		return result
	case fastjson.TypeString:
		return string(v.GetStringBytes())
	case fastjson.TypeNumber:
		// Try to parse as int64 first, then fall back to float64
		f := v.GetFloat64()
		if f == float64(int64(f)) {
			return int64(f)
		}
		return f
	case fastjson.TypeTrue:
		return true
	case fastjson.TypeFalse:
		return false
	case fastjson.TypeNull:
		return nil
	default:
		return nil
	}
}

// obscureGeneric recursively processes a generic structure and obscures known fields
func obscureGeneric(input interface{}) interface{} {
	switch v := input.(type) {
	case map[string]interface{}:
		return obscureMap(v)
	case []interface{}:
		return obscureArray(v)
	default:
		return input
	}
}

// obscureMap processes a map and obscures known fields
func obscureMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range m {
		if obscurableFields[key] {
			result[key] = obscureField(key, value)
		} else {
			// For unknown fields, recursively process if they're nested structures
			result[key] = obscureGeneric(value)
		}
	}

	return result
}

// obscureArray processes an array and obscures each element
func obscureArray(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, item := range arr {
		result[i] = obscureGeneric(item)
	}
	return result
}

// obscureField applies field-specific obscuration logic
func obscureField(fieldName string, value interface{}) interface{} {
	// Get the ID for deterministic hashing - we'll look for it in the parent context
	// For now, use empty string as fallback (this will be improved when we have context)
	id := ""

	switch fieldName {
	case "id":
		// ID is typically kept but could be hashed if needed
		if str, ok := value.(string); ok {
			return str
		}
	case "name":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicName(id, str)
		}
	case "first_name":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicFirstName(id, str)
		}
	case "last_name":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicLastName(id, str)
		}
	case "middle_name":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicMiddleName(id, str)
		}
	case "street":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicStreet(id, str)
		}
	case "email":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicEmail(id, str)
		}
	case "phone_number":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicPhone(id, str)
		}
	case "address":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicAddress(id, str)
		}
	case "city":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicCity(id, str)
		}
	case "state":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicState(id, str)
		}
	case "zip_code":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicZipCode(id, str)
		}
	case "county":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicCounty(id, str)
		}
	case "country":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicCountry(id, str)
		}
	case "tax_id":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicTaxID(id, str)
		}
	case "ssn":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicSSN(id, str)
		}
	case "date_of_birth":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicDateOfBirth(id, str)
		}
	case "gender":
		if str, ok := value.(string); ok {
			return data.GenerateDeterministicGender(id, str)
		}
	case "integer_value":
		if num, ok := toInt64(value); ok {
			return data.GenerateDeterministicInteger(id, num)
		}
	case "float_value":
		if fval, ok := toFloat64(value); ok {
			return data.GenerateDeterministicFloat(id, fval)
		}
	case "passport":
		return obscurePassport(value, id)
	case "driver_license":
		return obscureDriverLicense(value, id)
	case "bank_accounts":
		return obscureBankAccounts(value, id)
	}

	return value
}

// Helper functions for type conversion
func toInt64(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case float64:
		return int64(val), true
	case int64:
		return val, true
	case int:
		return int64(val), true
	default:
		return 0, false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int64:
		return float64(val), true
	case int:
		return float64(val), true
	default:
		return 0, false
	}
}

// Helper function to obscure passport data
func obscurePassport(value interface{}, id string) interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		result := make(map[string]interface{})
		for k, v := range m {
			if str, ok := v.(string); ok {
				result[k] = data.GenerateDeterministicPassportNumber(id+fmt.Sprintf("passport_%s_", k), str)
			} else {
				result[k] = v
			}
		}
		return result
	}
	return value
}

// Helper function to obscure driver license data
func obscureDriverLicense(value interface{}, id string) interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		result := make(map[string]interface{})
		for k, v := range m {
			if str, ok := v.(string); ok {
				result[k] = data.GenerateDeterministicDriverLicenseNumber(id+fmt.Sprintf("license_%s_", k), str)
			} else {
				result[k] = v
			}
		}
		return result
	}
	return value
}

// Helper function to obscure bank accounts
func obscureBankAccounts(value interface{}, id string) interface{} {
	if arr, ok := value.([]interface{}); ok {
		result := make([]interface{}, len(arr))
		for i, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				obscured := make(map[string]interface{})
				for k, v := range m {
					if str, ok := v.(string); ok {
						obscured[k] = data.GenerateDeterministicAccountName(id, str, i)
					} else {
						obscured[k] = v
					}
				}
				result[i] = obscured
			} else {
				result[i] = item
			}
		}
		return result
	}
	return value
}
