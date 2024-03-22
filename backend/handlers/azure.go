package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleAzure(logger *log.Logger, c *gin.Context, data []interface{}) {
	logger.Printf("Handling Azure request with data: %v", data)
	c.JSON(http.StatusOK, gin.H{"azure": data})
}
