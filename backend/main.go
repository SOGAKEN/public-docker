package main

import (
	"backend/handlers"
	"fmt"
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

		fmt.Println("kitayo")

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
