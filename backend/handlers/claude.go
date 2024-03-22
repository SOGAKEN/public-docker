package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
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
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Stop       string `json:"stop"`
}

func HandleClaude(c *gin.Context, data []interface{}) {
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty 'anthropic' data in the request"})
		return
	}

	req, ok := data[0].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Claude request format"})
		return
	}

	claudeReq := ClaudeRequest{}
	claudeReq.Model = req["model"].(string)

	messages, ok := req["messages"].([]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messages format"})
		return
	}

	if len(messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty messages in the request"})
		return
	}

	msgMap, ok := messages[0].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
		return
	}

	content, contentOk := msgMap["content"].(string)
	if !contentOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
		return
	}

	summaryPrompt := "この文章を要約してください。結果のみ返してください：" + content

	claudeReq.Messages = append(claudeReq.Messages, Message{Role: "Human", Content: summaryPrompt})

	// AWS SDKの設定
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// AWS Bedrockランタイムのクライアントを初期化
	brc := bedrockruntime.NewFromConfig(cfg)

	prompt := ""
	for _, msg := range claudeReq.Messages {
		prompt += "\n" + msg.Role + ": " + msg.Content
	}
	prompt += "\nAssistant:"

	payload := map[string]interface{}{
		"prompt":               prompt,
		"max_tokens_to_sample": 2048,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// タイムアウト設定付きのコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// メッセージをClaude APIに送信
	output, err := brc.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(claudeReq.Model),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var claudeResp ClaudeResponse
	err = json.Unmarshal(output.Body, &claudeResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling Claude response: " + err.Error(), "claude_response": string(output.Body)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": claudeReq.Model,
		"claude": gin.H{
			"content": claudeResp.Completion,
			"role":    "assistant",
		},
	})
}
