package iz2go

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var routes = map[string]*HandlerInfo{}

func NewRouter() *gin.Engine {
	router := gin.Default()
	fmt.Println("routes: ", routes)
	for path, route := range routes {
		router.Handle(route.Method, path, route.Handler)
	}
	return router
}

func RegisterRoute(path string, handler *HandlerInfo) {
	routes[path] = handler
}
