package middleware

import (
	"fmt"
	"net/http"
	"ss-catalog-service/internal/domain"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret would ideally come from config
var jwtSecret = []byte("your-256-bit-secret")

// AuthMiddleware extracts authentication claims from a JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		var userCtx domain.UserContext
		
		// Extract SellerID
		if sid, ok := claims["seller_id"].(float64); ok {
			id := int(sid)
			userCtx.SellerID = &id
		}

		// Extract Roles
		if roles, ok := claims["roles"].([]interface{}); ok {
			userCtx.Roles = make([]string, len(roles))
			for i, r := range roles {
				userCtx.Roles[i] = fmt.Sprint(r)
			}
		}

		if userCtx.SellerID != nil {
			newCtx := domain.ContextWithUser(c.Request.Context(), userCtx)
			c.Request = c.Request.WithContext(newCtx)
			c.Set("seller_id", *userCtx.SellerID)
		}

		c.Next()
	}
}

// RequireAuth is a helper middleware to block unauthorized requests.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := domain.UserFromContext(c.Request.Context())
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		c.Next()
	}
}
