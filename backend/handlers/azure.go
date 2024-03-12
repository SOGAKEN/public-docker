package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleAzure(c *gin.Context, data []interface{}) {
	// Azure APIへのリクエストを処理するロジックを実装する
	c.JSON(http.StatusOK, gin.H{"azure": data})
}
