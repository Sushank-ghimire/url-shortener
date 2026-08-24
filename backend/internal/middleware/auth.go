package middleware

import (
	"net/http"
	"strings"

	"github.com/Sushank-ghimire/url-shortener/internal/auth"
	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func AuthRequired(
	jwtService *auth.JWTService,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "authorization header required",
				},
			)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid authorization header",
				},
			)
			return
		}

		token := parts[1]
		userId, err := jwtService.ValidateAccessToken(token)

		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid or expired access token",
				},
			)
			return
		}
		ctx.Set(UserIDKey, userId)
		ctx.Next()
	}
}
