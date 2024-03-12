package handlers

import (
	"github.com/gin-gonic/gin"
)

// APIHandler defines the interface for handling requests to different APIs.
type APIHandler interface {
	HandleRequest(c *gin.Context, requestData interface{}) (interface{}, error)
}
