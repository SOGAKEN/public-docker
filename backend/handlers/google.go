package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGoogle(c *gin.Context, data []interface{}) {
	// Google APIへのリクエストを処理するロジックを実装する
	c.JSON(http.StatusOK, gin.H{"google": data})
}
