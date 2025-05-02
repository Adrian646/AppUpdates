package handler

import (
	"fmt"
	"github.com/Adrian646/AppUpdates/backend/internal/feeds/ios"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IOSFeedHandler(c *gin.Context) {
	appID := c.Param("appID")

	feed, err := ios.GetCurrentAppData(appID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feed)
}

func AndroidFeedHandler(c *gin.Context) {
	appID := c.Param("appID")
	fmt.Println("Handling android: ", appID)
	c.JSON(200, gin.H{})
}
