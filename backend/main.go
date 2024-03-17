package main

import (
	"backend/handlers"
	"fmt"
	"io"
	"net/http"
	"os"

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
	orgin := os.Getenv("AUTH_URL")
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{orgin}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	username := os.Getenv("BASIC_AUTH_USER")
	password := os.Getenv("BASIC_AUTH_PASS")

	fmt.Println(username, password)

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	authorized.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "認証済み"})
	})

	r.POST("/api/login", handlers.Login)

	r.POST("/api/summary", func(c *gin.Context) {
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for key, value := range body.Data {
			switch key {
			case "openai":
				handlers.HandleOpenAI(c, value)
			case "google":
				handlers.HandleGoogle(c, value)
			case "azure":
				handlers.HandleAzure(c, value)
			case "anthropic":
				handlers.HandleClaude(c, value)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key"})
				return
			}
		}
	})

	r.Run(":8080")
}
