package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func IOSFeedHandler(c *gin.Context) {
	appID := c.Param("appID")
	fmt.Println("Handling ios: ", appID)
	c.JSON(200, gin.H{})
}

func AndroidFeedHandler(c *gin.Context) {
	appID := c.Param("appID")
	fmt.Println("Handling android: ", appID)
	c.JSON(200, gin.H{})
}
