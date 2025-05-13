package iz2go

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

var routes = map[string]*HandlerInfo{}

type Engine struct {
	*gin.Engine
}

type SwaggerRenderConfig struct {
	SwaggerPath string
	OpenApiPath string
	Info        *Info
}

func (e *Engine) RenderSwagger(config *SwaggerRenderConfig) {
	if config.SwaggerPath == "" {
		config.SwaggerPath = "/docs"
	}
	if config.OpenApiPath == "" {
		config.OpenApiPath = "/openapi.json"
	}
	if config.Info == nil {
		config.Info = &Info{}
	}
	if config.Info.Title == "" {
		config.Info.Title = "API Documentation"
	}
	if config.Info.Description == "" {
		config.Info.Description = "API Documentation for the project"
	}
	if config.Info.Version == "" {
		config.Info.Version = "1.0.0"
	}

	swaggerConfig := GenerateSwagger(config.Info)
	e.GET(config.OpenApiPath, func(c *gin.Context) {
		c.JSON(200, swaggerConfig)
	})

	e.GET(config.SwaggerPath, func(c *gin.Context) {
		c.HTML(200, "swagger", gin.H{
			"info":        config.Info,
			"openApiPath": config.OpenApiPath,
		})
	})
}

func Default() *Engine {
	router := gin.Default()
	tmpl := template.Must(template.New("swagger").Parse(swaggerHTML))
	router.SetHTMLTemplate(tmpl)
	for path, route := range routes {
		router.Handle(route.Method, path, route.Handler)
	}
	return &Engine{
		Engine: router,
	}
}

func RegisterRoute(path string, handler *HandlerInfo) {
	routes[path] = handler
}
