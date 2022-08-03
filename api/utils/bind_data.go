package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindData(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
