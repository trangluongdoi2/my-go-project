package middleware

import (
	"go-backend-project/internal/auth"
	"go-backend-project/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.AbortWithError(c, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.AbortWithError(c, http.StatusUnauthorized, "Invalid authorization header format. Use: Bearer <token>")
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			utils.AbortWithError(c, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		if claims.TokenType != auth.AccessToken {
			utils.AbortWithError(c, http.StatusUnauthorized, "Invalid token type. Access token required")
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.AbortWithError(c, http.StatusUnauthorized, "User role not found")
			return
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		utils.AbortWithError(c, http.StatusForbidden, "Access denied. Insufficient permissions")
	}
}
