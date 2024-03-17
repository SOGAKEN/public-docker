package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/vertexai/genai"
	"github.com/gin-gonic/gin"
)

type VertexAIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []string `json:"Parts"`
		} `json:"Content"`
	} `json:"Candidates"`
}

func HandleGoogle(c *gin.Context, data []interface{}) {
	if len(data) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided"})
		return
	}

	firstData, ok := data[0].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	modelName, modelOk := firstData["model"].(string)
	messages, messagesOk := firstData["messages"].([]interface{})
	if !modelOk || !messagesOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data structure"})
		return
	}
	if len(messages) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided"})
		return
	}

	message, ok := messages[len(messages)-1].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
		return
	}

	content, contentOk := message["content"].(string)
	if !contentOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message content"})
		return
	}

	projectId := os.Getenv("PROJECT_ID")
	region := os.Getenv("REGION")

	parts, err := makeChatRequests(projectId, region, modelName, content)
	if err != nil {
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

func makeChatRequests(projectId, region, modelName, content string) ([]genai.Part, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectId, region)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.9)
	resp, err := model.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		return nil, err
	}

	var parts []genai.Part
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			fmt.Println(part)
			parts = append(parts, part)

		}
	}

	return parts, nil
}
