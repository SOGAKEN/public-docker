package main

import (
	"backend/handlers"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type RequestBody struct {
	Data map[string][]interface{} `json:"data"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// NextJSアプリケーションへのリバースプロキシを設定
	nextURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(nextURL)

	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	username := os.Getenv("BASIC_AUTH_USER")
	password := os.Getenv("BASIC_AUTH_PASS")

	fmt.Println(username, password)

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	authorized.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "認証済み"})
	})

	r.POST("/api/login", func(c *gin.Context) {
		var loginReq LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username := os.Getenv("BASIC_AUTH_USER")
		password := os.Getenv("BASIC_AUTH_PASS")

		fmt.Println(loginReq.Username)
		fmt.Println(loginReq.Password)

		if loginReq.Username == username && loginReq.Password == password {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": username,
				"exp":      time.Now().Add(time.Hour * 24).Unix(),
			})

			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.SetCookie("token", tokenString, 86400, "/", "localhost", false, true)
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		}
	})

	r.POST("/api", func(c *gin.Context) {
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
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key"})
				return
			}
		}
	})

	// NextJSアプリケーションへのリバースプロキシを設定
	r.NoRoute(func(c *gin.Context) {
		target, _ := url.Parse("http://localhost:8080")
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	r.Run(":8081")
}
