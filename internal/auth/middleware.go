package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware creates a middleware that validates JWT tokens
func JWTMiddleware(pkm *PublicKeyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := pkm.ValidateToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Store claims in context for downstream handlers
		c.Set("jwt_claims", claims)
		c.Next()
	}
}

// GetClaims retrieves JWT claims from the context
func GetClaims(c *gin.Context) (map[string]any, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}

	claimsMap, ok := claims.(map[string]any)
	return claimsMap, ok
}

// GetClaimString retrieves a string claim from the context
func GetClaimString(c *gin.Context, key string) (string, bool) {
	claims, exists := GetClaims(c)
	if !exists {
		return "", false
	}

	value, ok := claims[key]
	if !ok {
		return "", false
	}

	str, ok := value.(string)
	return str, ok
}
