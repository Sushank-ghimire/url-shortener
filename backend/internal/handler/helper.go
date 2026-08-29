package handler

import (
	"github.com/Sushank-ghimire/url-shortener/internal/middleware"
	"github.com/gin-gonic/gin"
)

func getUserID(ctx *gin.Context) (int64, bool) {
	value, exists := ctx.Get(middleware.UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}
