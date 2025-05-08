package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Method = "GET"

type Login struct{}

func (h *Login) Execute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
