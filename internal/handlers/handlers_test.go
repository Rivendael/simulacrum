package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simulacrum/internal/data"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/obscure", HandleObscure)
	return r
}

func TestHandleObscureSuccess(t *testing.T) {
	router := setupRouter()

	input := data.PersonalData{
		ID:          "user123",
		Name:        "John Doe",
		Email:       "john@example.com",
		Address:     "123 Main St",
		PhoneNumber: "555-1234",
		TaxID:       "123-45-6789",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response data.PersonalData
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, input.ID, response.ID)
	assert.NotEqual(t, input.Name, response.Name)
	assert.NotEqual(t, input.Email, response.Email)
	assert.NotEqual(t, input.TaxID, response.TaxID)
}

func TestHandleObscureMissingID(t *testing.T) {
	router := setupRouter()

	input := data.PersonalData{Name: "John Doe"}
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleObscureDeterministic(t *testing.T) {
	router := setupRouter()

	input := data.PersonalData{
		ID:    "user123",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	body, _ := json.Marshal(input)

	req1, _ := http.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	var response1 data.PersonalData
	json.Unmarshal(w1.Body.Bytes(), &response1)

	req2, _ := http.NewRequest("POST", "/obscure", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var response2 data.PersonalData
	json.Unmarshal(w2.Body.Bytes(), &response2)

	assert.Equal(t, response1.Name, response2.Name)
	assert.Equal(t, response1.Email, response2.Email)
	assert.Equal(t, response1.TaxID, response2.TaxID)
}
