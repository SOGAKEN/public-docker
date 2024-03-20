package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIHandler defines the interface for handling requests to different APIs.
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

func HandleSummary(c *gin.Context) {
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, value := range body.Data {
		switch key {
		case "openai":
			HandleOpenAI(c, value)
		case "google":
			HandleGoogle(c, value)
		case "azure":
			HandleAzure(c, value)
		case "anthropic":
			HandleClaude(c, value)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key"})
			return
		}
	}
}
