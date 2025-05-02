package main

import (
	"fmt"
	"github.com/Adrian646/AppUpdates/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../backend.env")

	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	r := gin.Default()

	r.GET("/api/v1/feeds/android/:appID", handler.AndroidFeedHandler)
	r.GET("/api/v1/feeds/ios/:appID", handler.IOSFeedHandler)

	err = r.Run()

	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
}
