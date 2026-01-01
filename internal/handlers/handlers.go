package handlers

import (
	"net/http"

	"simulacrum/internal/data"

	"github.com/gin-gonic/gin"
)

func HandleObscure(c *gin.Context) {
	var input data.PersonalData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if input.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	obscured := data.ObscureData(input)
	c.JSON(http.StatusOK, obscured)
}
