package handler

import (
	"github.com/Sushank-ghimire/url-shortener/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

func buildAnalytics(
	c *gin.Context,
	urlID int64,
) *model.Analytics {
	userAgent := c.GetHeader("User-Agent")
	ua := useragent.New(userAgent)

	browser, _ := ua.Browser()

	analytics := &model.Analytics{
		URLID:     urlID,
		IpAddress: c.ClientIP(),
		Browser:   browser,
		Os:        ua.OS(),
		Device:    ua.Model(),
		Referer:   c.GetHeader("Referer"),
	}
	return analytics
}
