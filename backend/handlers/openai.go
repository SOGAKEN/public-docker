package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
	RequestModel string `json:"request_model"`
}

func HandleOpenAI(logger *log.Logger, c *gin.Context, data []interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(data))

	for _, v := range data {
		go func(v interface{}) {
			defer wg.Done()

			newC := gin.Context{}
			newC.Request = c.Request.Clone(c.Request.Context())
			newC.Writer = &ResponseWriterMock{ResponseWriter: c.Writer, body: &bytes.Buffer{}}

			req, ok := v.(map[string]interface{})
			if !ok {
				logger.Printf("Invalid OpenAI request: %v", v)
				newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OpenAI request"})
				return
			}

			openaiReq := OpenAIRequest{}
			openaiReq.Model = req["model"].(string)

			messages, ok := req["messages"].([]interface{})
			if !ok {
				logger.Printf("Invalid messages format: %v", req["messages"])
				newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messages format"})
				return
			}

			for _, msg := range messages {
				msgMap, ok := msg.(map[string]interface{})
				if !ok {
					logger.Printf("Invalid message format: %v", msg)
					newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
					return
				}

				role, roleOk := msgMap["role"].(string)
				content, contentOk := msgMap["content"].(string)
				if !roleOk || !contentOk {
					logger.Printf("Invalid message format: %v", msgMap)
					newC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
					return
				}

				openaiReq.Messages = append(openaiReq.Messages, Message{Role: role, Content: content})
			}

			jsonReq, err := json.Marshal(openaiReq)
			if err != nil {
				logger.Printf("Error creating OpenAI request: %v", err)
				newC.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			apiKey := os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				logger.Printf("OPENAI_API_KEY environment variable not set")
				newC.JSON(http.StatusInternalServerError, gin.H{"error": "OPENAI_API_KEY environment variable not set"})
				return
			}

			httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonReq))
			if err != nil {
				logger.Printf("Error creating OpenAI request: %v", err)
				newC.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("Authorization", "Bearer "+apiKey)

			client := &http.Client{}
			resp, err := client.Do(httpReq)
			if err != nil {
				logger.Printf("Error sending OpenAI request: %v", err)
				newC.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Printf("Error reading OpenAI response body: %v", err)
				newC.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var openaiResp OpenAIResponse
			openaiResp.RequestModel = openaiReq.Model
			err = json.Unmarshal(body, &openaiResp)
			if err != nil {
				logger.Printf("Error unmarshaling OpenAI response: %v", err)
				newC.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling OpenAI response: " + err.Error(), "openai_response": string(body)})
				return
			}

			if openaiResp.Error.Message != "" {
				logger.Printf("OpenAI API error: %v", openaiResp.Error)
				newC.JSON(http.StatusInternalServerError, gin.H{"error": openaiResp.Error.Message, "openai_response": openaiResp})
				return
			}

			if len(openaiResp.Choices) == 0 {
				logger.Printf("No choices in OpenAI response: %v", openaiResp)
				newC.JSON(http.StatusInternalServerError, gin.H{"error": "No choices in OpenAI response", "openai_response": openaiResp})
				return
			}

			newC.JSON(http.StatusOK, gin.H{
				"openai": openaiResp.Choices[0].Message,
				"model":  openaiResp.RequestModel,
			})
		}(v)
	}

	wg.Wait()
}

type ResponseWriterMock struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *ResponseWriterMock) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *ResponseWriterMock) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
