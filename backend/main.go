package main

import (
	"backend/handlers"
	"backend/middleware"
	"fmt"
	"io"
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
		// 環境変数 GOOGLE_APPLICATION_CREDENTIALS を一時ファイルに設定
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", tmpfile.Name())
	}

	r := gin.Default()
	origins := strings.Split(os.Getenv("AUTH_URL"), ",")
	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// username := os.Getenv("BASIC_AUTH_USER")
	// password := os.Getenv("BASIC_AUTH_PASS")
	//
	// fmt.Println(username, password)
	//
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	username: password,
	// }))
	//
	// authorized.POST("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"message": "認証済み"})
	// })

	r.POST("/api/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/api/summary", handlers.HandleSummary)

	r.Run(":8080")
}
