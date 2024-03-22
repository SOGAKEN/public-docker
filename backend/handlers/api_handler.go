package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIHandler interface {
	HandleRequest(c *gin.Context, requestData interface{}) (interface{}, error)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Data map[string][]interface{} `json:"data"`
}

func HandleSummary(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for key, value := range body.Data {
			switch key {
			case "openai":
				HandleOpenAI(logger, c, value)
			case "google":
				HandleGoogle(logger, c, value)
			case "azure":
				HandleAzure(logger, c, value)
			case "anthropic":
				HandleClaude(logger, c, value)
			default:
				logger.Printf("Invalid API key: %s", key)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key"})
				return
			}
		}
	}
}
