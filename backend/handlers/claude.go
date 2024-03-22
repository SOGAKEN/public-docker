package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gin-gonic/gin"
)

type ClaudeRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ClaudeResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Role       string    `json:"role"`
	Content    []Content `json:"content"`
	Model      string    `json:"model"`
	StopReason string    `json:"stop_reason"`
	Usage      Usage     `json:"usage"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func HandleClaude(logger *log.Logger, c *gin.Context, data []interface{}) {
	if len(data) == 0 {
		logger.Printf("Empty 'anthropic' data in the request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty 'anthropic' data in the request"})
		return
	}

	req, ok := data[0].(map[string]interface{})
	if !ok {
		logger.Printf("Invalid Claude request format: %v", data[0])
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Claude request format"})
		return
	}

	claudeReq := ClaudeRequest{}
	claudeReq.Model = req["model"].(string)

	messages, ok := req["messages"].([]interface{})
	if !ok {
		logger.Printf("Invalid messages format: %v", req["messages"])
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messages format"})
		return
	}

	if len(messages) == 0 {
		logger.Printf("Empty messages in the request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty messages in the request"})
		return
	}

	var contentParts []string
	for _, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			logger.Printf("Invalid message format: %v", msg)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
			return
		}

		content, contentOk := msgMap["content"].(string)
		if !contentOk {
			logger.Printf("Invalid message format: %v", msgMap)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
			return
		}

		contentParts = append(contentParts, content)
	}

	combinedContent := strings.Join(contentParts, " ")

	claudeReq.Messages = append(claudeReq.Messages, Message{Role: "user", Content: combinedContent})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		logger.Printf("Error loading AWS config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	brc := bedrockruntime.NewFromConfig(cfg)

	anthropicVersion := os.Getenv("ANTHROPIC_VERSION")
	if anthropicVersion == "" {
		anthropicVersion = "bedrock-2023-05-31"
	}

	message := bedrockruntime.InvokeModelInput{
		ModelId: aws.String("anthropic." + claudeReq.Model),
		Body: func() []byte {
			payload := map[string]interface{}{
				"max_tokens":        2048,
				"messages":          claudeReq.Messages,
				"anthropic_version": anthropicVersion,
			}
			payloadBytes, _ := json.Marshal(payload)
			return payloadBytes
		}(),
		ContentType: aws.String("application/json"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	output, err := brc.InvokeModel(ctx, &message)
	if err != nil {
		logger.Printf("Error invoking Claude model: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var claudeResp ClaudeResponse
	err = json.Unmarshal(output.Body, &claudeResp)
	if err != nil {
		logger.Printf("Error unmarshaling Claude response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling Claude response: " + err.Error(), "claude_response": string(output.Body)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": claudeReq.Model,
		"claude": gin.H{
			"content": func() string {
				if len(claudeResp.Content) > 0 {
					return claudeResp.Content[0].Text
				}
				return ""
			}(),
			"role": "assistant",
		},
	})
}
