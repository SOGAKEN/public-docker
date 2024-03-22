package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/vertexai/genai"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type VertexAIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []string `json:"Parts"`
		} `json:"Content"`
	} `json:"Candidates"`
}

func HandleGoogle(logger *log.Logger, c *gin.Context, data []interface{}) {
	if len(data) < 1 {
		logger.Printf("No data provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided"})
		return
	}

	firstData, ok := data[0].(map[string]interface{})
	if !ok {
		logger.Printf("Invalid data format: %v", data[0])
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	modelName, modelOk := firstData["model"].(string)
	messages, messagesOk := firstData["messages"].([]interface{})
	if !modelOk || !messagesOk {
		logger.Printf("Invalid data structure: %v", firstData)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data structure"})
		return
	}
	if len(messages) < 1 {
		logger.Printf("No messages provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided"})
		return
	}

	message, ok := messages[len(messages)-1].(map[string]interface{})
	if !ok {
		logger.Printf("Invalid message format: %v", messages[len(messages)-1])
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
		return
	}

	content, contentOk := message["content"].(string)
	if !contentOk {
		logger.Printf("Invalid message content: %v", message)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message content"})
		return
	}

	projectId := os.Getenv("PROJECT_ID")
	region := os.Getenv("REGION")

	parts, err := makeChatRequests(logger, projectId, region, modelName, content)
	if err != nil {
		logger.Printf("Error making chat requests: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": modelName,
		"google": gin.H{
			"content": parts[0],
			"role":    "assistant",
		},
	})
}

func makeChatRequests(logger *log.Logger, projectId, region, modelName, content string) ([]genai.Part, error) {
	ctx := context.Background()

	credsJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
	var opts []option.ClientOption
	if credsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credsJSON)))
	}

	client, err := genai.NewClient(ctx, projectId, region, opts...)
	if err != nil {
		logger.Printf("Error creating client: %v", err)
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.9)
	resp, err := model.GenerateContent(ctx, genai.Text("Provide a summary for the following article: "+content))
	if err != nil {
		logger.Printf("Error generating content: %v", err)
		return nil, fmt.Errorf("error generating content: %v", err)
	}

	var parts []genai.Part
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			parts = append(parts, part)
		}
	}

	return parts, nil
}
