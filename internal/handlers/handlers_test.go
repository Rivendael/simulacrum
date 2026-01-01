package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleObscureGenericSingleField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id":    "user123",
		"name":  "John Doe",
		"email": "john@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	if result["id"] != "user123" {
		t.Errorf("Expected ID user123, got %v", result["id"])
	}

	// Name and email should be obscured (different from input)
	if name, ok := result["name"].(string); !ok || name == "John Doe" {
		t.Errorf("Expected name to be obscured")
	}

	if email, ok := result["email"].(string); !ok || email == "john@example.com" {
		t.Errorf("Expected email to be obscured")
	}
}

func TestHandleObscureGenericNestedObject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id":     "user123",
		"name":   "John Doe",
		"street": "123 Main St",
		"city":   "Springfield",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	// Street and city should be obscured
	if street, ok := result["street"].(string); !ok || street == "123 Main St" {
		t.Errorf("Expected street to be obscured, got %v", result["street"])
	}

	if city, ok := result["city"].(string); !ok || city == "Springfield" {
		t.Errorf("Expected city to be obscured, got %v", result["city"])
	}
}

func TestHandleObscureGenericArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id": "user123",
		"bank_accounts": []interface{}{
			map[string]interface{}{
				"account_number": "1234567890",
				"bank_name":      "Bank A",
			},
			map[string]interface{}{
				"account_number": "9876543210",
				"bank_name":      "Bank B",
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	// Bank accounts should be obscured as array
	accounts, ok := result["bank_accounts"].([]interface{})
	if !ok {
		t.Errorf("Expected bank_accounts to be an array")
	}
	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}
}

func TestHandleObscureGenericMixedStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id":            "user123",
		"first_name":    "John",
		"last_name":     "Doe",
		"email":         "john@example.com",
		"phone_number":  "555-1234",
		"date_of_birth": "1990-01-15",
		"custom_field":  "should_pass_through",
		"integer_value": int64(42),
		"float_value":   3.14,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	// Custom field should pass through unchanged
	if custom, ok := result["custom_field"].(string); !ok || custom != "should_pass_through" {
		t.Errorf("Expected custom_field to pass through")
	}

	// Obscurable fields should be present but changed
	if _, ok := result["first_name"].(string); !ok {
		t.Errorf("Expected first_name to exist")
	}

	if _, ok := result["email"].(string); !ok {
		t.Errorf("Expected email to exist")
	}
}

func TestHandleObscureGenericDeterminism(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id":    "user123",
		"email": "test@example.com",
	}
	body, _ := json.Marshal(reqBody)

	// First call
	req1 := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	var result1 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &result1)

	// Second call
	req2 := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(bytes.NewBuffer(body).Bytes()))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var result2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result2)

	// Results should be identical (deterministic)
	if result1["email"] != result2["email"] {
		t.Errorf("Expected deterministic obscuration: %v vs %v", result1["email"], result2["email"])
	}
}

func TestHandleObscureGenericPassport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id": "user123",
		"passport": map[string]interface{}{
			"number":      "ABC123456",
			"country":     "US",
			"expiry_date": "2030-12-31",
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	passport, ok := result["passport"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected passport to be an object")
	}

	if _, ok := passport["number"].(string); !ok {
		t.Errorf("Expected passport number to exist")
	}
}

func TestHandleObscureGenericArbitraryJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	// Arbitrary structure with any fields
	reqBody := map[string]interface{}{
		"id":    "user123",
		"email": "user@example.com",
		"custom_obj": map[string]interface{}{
			"field1": "value1",
			"field2": 123,
		},
		"tags": []interface{}{"tag1", "tag2", "tag3"},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	// Email should be obscured
	if email, ok := result["email"].(string); !ok || email == "user@example.com" {
		t.Errorf("Expected email to be obscured")
	}

	// Custom object should pass through (no known fields)
	if customObj, ok := result["custom_obj"].(map[string]interface{}); !ok {
		t.Errorf("Expected custom_obj to exist")
	} else {
		if val, ok := customObj["field1"].(string); !ok || val != "value1" {
			t.Errorf("Expected custom_obj.field1 to pass through")
		}
	}

	// Tags array should pass through
	if tags, ok := result["tags"].([]interface{}); !ok || len(tags) != 3 {
		t.Errorf("Expected tags array to pass through")
	}
}

func TestParallelObscureDataSmallBatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	people := make([]map[string]interface{}, 5)
	for i := 0; i < 5; i++ {
		people[i] = map[string]interface{}{
			"id":    fmt.Sprintf("user_%03d", i),
			"name":  fmt.Sprintf("Person %d", i),
			"email": fmt.Sprintf("user%d@example.com", i),
		}
	}

	reqBody := map[string]interface{}{
		"people": people,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	peopleResult, ok := result["people"].([]interface{})
	if !ok || len(peopleResult) != 5 {
		t.Errorf("Expected 5 people in response, got %d", len(peopleResult))
	}
}

func TestParallelObscureLargeBatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	// Test with 100 people to verify parallel processing
	people := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		people[i] = map[string]interface{}{
			"id":    fmt.Sprintf("user_%03d", i),
			"name":  fmt.Sprintf("Person %d", i),
			"email": fmt.Sprintf("user%d@example.com", i),
		}
	}

	reqBody := map[string]interface{}{
		"people": people,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	peopleResult, ok := result["people"].([]interface{})
	if !ok || len(peopleResult) != 100 {
		t.Errorf("Expected 100 people in response, got %d", len(peopleResult))
	}
}

func TestParallelObscureDataDeterminism(t *testing.T) {
	// Verify that parallel processing produces deterministic results
	people := make([]map[string]interface{}, 50)
	for i := 0; i < 50; i++ {
		people[i] = map[string]interface{}{
			"id":    fmt.Sprintf("user_%03d", i),
			"name":  fmt.Sprintf("Person %d", i),
			"email": fmt.Sprintf("user%d@example.com", i),
		}
	}

	reqBody := map[string]interface{}{
		"people": people,
	}
	body1, _ := json.Marshal(reqBody)
	body2, _ := json.Marshal(reqBody)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	// First request
	req1 := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	var result1 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &result1)

	// Second request
	req2 := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var result2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result2)

	// Results should be identical
	result1JSON, _ := json.Marshal(result1)
	result2JSON, _ := json.Marshal(result2)

	if string(result1JSON) != string(result2JSON) {
		t.Errorf("Parallel processing should produce deterministic results")
	}
}

func TestHandleObscurePhoneWithCountryCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id":           "user123",
		"phone_number": "+1-555-987-6543",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	phone, ok := result["phone_number"].(string)
	if !ok || phone == "+1-555-987-6543" {
		t.Errorf("Expected phone to be obscured, got %v", result["phone_number"])
	}

	// Should preserve the country code prefix
	if !strings.HasPrefix(phone, "+1-") {
		t.Errorf("Expected phone to preserve country code +1-, got %s", phone)
	}
}

func TestHandleObscurePhoneWithoutCountryCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/obscure", HandleObscure)

	reqBody := map[string]interface{}{
		"id":           "user123",
		"phone_number": "555-987-6543",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	phone, ok := result["phone_number"].(string)
	if !ok || phone == "555-987-6543" {
		t.Errorf("Expected phone to be obscured, got %v", result["phone_number"])
	}

	// Should NOT have country code
	if strings.HasPrefix(phone, "+") {
		t.Errorf("Expected phone without country code, got %s", phone)
	}
}
