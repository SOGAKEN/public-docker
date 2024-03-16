package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

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
}

func HandleClaude(c *gin.Context, data []interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(data))

	for _, v := range data {
		go func(v interface{}) {
			defer wg.Done()

			newC := gin.Context{}
			newC.Request = c.Request.Clone(c.Request.Context())

			req, ok := v.(map[string]interface{})
			if !ok {
				newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Claude request"})
				return
			}

			claudeReq := ClaudeRequest{}
			claudeReq.Model = req["model"].(string)

			messages, ok := req["messages"].([]interface{})
			if !ok {
				newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messages format"})
				return
			}

			for _, msg := range messages {
				msgMap, ok := msg.(map[string]interface{})
				if !ok {
					newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
					return
				}

				role, roleOk := msgMap["role"].(string)
				content, contentOk := msgMap["content"].(string)
				if !roleOk || !contentOk {
					newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
					return
				}

				claudeReq.Messages = append(claudeReq.Messages, Message{Role: role, Content: content})
			}

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
				newC.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// AWS Bedrockランタイムのクライアントを初期化
			brc := bedrockruntime.NewFromConfig(cfg)

			prompt := ""
			for _, msg := range claudeReq.Messages {
				prompt += "\n" + msg.Role + ": " + msg.Content
			}
			prompt += "\nAssistant:"

			claudeInput := &bedrockruntime.InvokeModelInput{
				Body:        []byte(`{"prompt":"` + prompt + `","max_tokens_to_sample":2048}`),
				ModelId:     aws.String(claudeReq.Model),
				ContentType: aws.String("application/json"),
			}

			// メッセージをClaude APIに送信
			output, err := brc.InvokeModel(context.Background(), claudeInput)
			if err != nil {
				newC.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var claudeResp ClaudeResponse
			err = json.Unmarshal(output.Body, &claudeResp)
			if err != nil {
				newC.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling Claude response: " + err.Error(), "claude_response": string(output.Body)})
				return
			}

			newC.JSON(http.StatusOK, gin.H{
				"claude": gin.H{
					"content": claudeResp.Completion,
				},
				"model": claudeReq.Model,
			})
		}(v)
	}

	wg.Wait()
}
