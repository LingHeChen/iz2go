package iz2go

import "github.com/gin-gonic/gin"

type HandlerInfo struct {
	Method  string
	Handler gin.HandlerFunc
}

type Post struct{}

func (p *Post) GetMethod() string {
	return "POST"
}

type Get struct{}

func (g *Get) GetMethod() string {
	return "GET"
}
