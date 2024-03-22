package main

import (
	"backend/handlers"
	"backend/middleware"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type RequestBody struct {
	Data map[string][]interface{} `json:"data"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
	if creds != "" {
		tmpfile, err := os.CreateTemp("", "gcp-creds-*.json")
		if err != nil {
			fmt.Println("Failed to create temp file for credentials", err)
			return
		}
		if _, err := io.WriteString(tmpfile, creds); err != nil {
			fmt.Println("Failed to write to temp credentials file", err)
			return
		}
		if err := tmpfile.Close(); err != nil {
			fmt.Println("Failed to close temp credentials file", err)
			return
		}
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", tmpfile.Name())
	}

	logFilePath := "./application.log"

	// ロガーの設定
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	r := gin.Default()
	origins := strings.Split(os.Getenv("AUTH_URL"), ",")
	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	r.POST("/api/login", handlers.LoginHandler(logger))

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(logger))
	protected.POST("/api/summary", handlers.HandleSummary(logger))

	r.Run(":8080")
}
