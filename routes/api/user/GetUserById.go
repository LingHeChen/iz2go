package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Method = "GET"

type GetUserById struct {
}

func (h *GetUserById) Execute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
