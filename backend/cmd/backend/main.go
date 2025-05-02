package main

import (
	"fmt"
	"github.com/Adrian646/AppUpdates/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load("../backend.env")

	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	r := gin.Default()

	r.Use(checkToken)

	r.GET("/api/feeds/android/:appID", handler.AndroidFeedHandler)
	r.GET("/api/feeds/ios/:appID", handler.IOSFeedHandler)

	err = r.Run()

	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
}

func checkToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != os.Getenv("API_KEY") {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
}
