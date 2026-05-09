package middleware

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"ss-catalog-service/config"
	"ss-catalog-service/internal/domain"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware extracts authentication claims from a JWT using RS256.
func AuthMiddleware(cfg config.JWTConfig) gin.HandlerFunc {
	// Pre-load public key for performance
	verifyKey, err := loadPublicKey(cfg.PublicKeyPath)
	if err != nil {
		fmt.Printf("Warning: failed to load JWT public key: %v\n", err)
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return verifyKey, nil
		}, 
		jwt.WithIssuer(cfg.Issuer),
		jwt.WithAudience(cfg.Audience),
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		var userCtx domain.UserContext
		
		// Extract UserID (sub)
		if sub, ok := claims["sub"].(string); ok {
			userCtx.UserID = sub
		}

		// Extract SellerID (custom claim)
		if sid, ok := claims["seller_id"].(float64); ok {
			id := int(sid)
			userCtx.SellerID = &id
		}

		// Extract Roles (Enterprise mapping)
		if roles, ok := claims["roles"].([]interface{}); ok {
			userCtx.Roles = make([]string, len(roles))
			for i, r := range roles {
				userCtx.Roles[i] = fmt.Sprint(r)
			}
		}

		// Set context
		newCtx := domain.ContextWithUser(c.Request.Context(), userCtx)
		c.Request = c.Request.WithContext(newCtx)
		
		c.Next()
	}
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(keyData)
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
