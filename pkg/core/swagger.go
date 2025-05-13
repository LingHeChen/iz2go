package iz2go

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

const swaggerHTML = `<!DOCTYPE html>
<html>
<head>
    <title>{{ .info.Title }}</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	{{- if .icon }}
	<link rel="icon" href="{{ .icon }}" />
	{{- end }}
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "{{ .openApiPath }}",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ]
            });
        }
    </script>
</body>
</html>`

// SwaggerConfig 表示 Swagger 配置
type SwaggerConfig struct {
	Swagger     string                `json:"swagger"`
	Info        Info                  `json:"info"`
	Paths       map[string]PathItem   `json:"paths"`
	Definitions map[string]Definition `json:"definitions"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

type Operation struct {
	Tags       []string            `json:"tags"`
	Summary    string              `json:"summary"`
	Parameters []Parameter         `json:"parameters"`
	Responses  map[string]Response `json:"responses"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Required    bool    `json:"required"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Schema      *Schema `json:"schema,omitempty"`
}

type Response struct {
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

type Schema struct {
	Type        string              `json:"type"`
	Format      string              `json:"format"`
	Description string              `json:"description"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Items       *Schema             `json:"items,omitempty"`
	Ref         string              `json:"$ref,omitempty"`
}

type Property struct {
	Type        string              `json:"type"`
	Description string              `json:"description,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Items       *Schema             `json:"items,omitempty"`
}

type Definition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// GenerateSwagger 生成 Swagger 配置
func GenerateSwagger(info *Info) *SwaggerConfig {
	config := &SwaggerConfig{
		Swagger:     "2.0",
		Info:        *info,
		Paths:       make(map[string]PathItem),
		Definitions: make(map[string]Definition),
	}

	for path, handler := range routes {
		if handler.ApiName == "" {
			pathParts := strings.ReplaceAll(path, "/", "_")
			handler.ApiName = pathParts
		}
		handlerType := reflect.TypeOf(handler.Handler)
		if handlerType.Kind() == reflect.Ptr {
			handlerType = handlerType.Elem()
		}

		method := getMethodFromHandler(handler)

		// 获取请求参数类型
		requestType := getRequestType(handler)
		if requestType != nil {
			// 生成参数定义
			parameters, definitions := generateParameters(requestType, handler)
			for name, definition := range definitions {
				config.Definitions[name] = definition
			}

			// 生成响应定义
			responses, responsesDefinitions := generateResponses(handler)
			for name, definition := range responsesDefinitions {
				config.Definitions[name] = definition
			}

			// 创建操作
			operation := &Operation{
				Tags:       []string{handlerType.Name()},
				Summary:    getSummaryFromHandler(handler),
				Parameters: parameters,
				Responses:  responses,
			}

			// 添加到路径
			pathItem := config.Paths[path]
			switch method {
			case "GET":
				pathItem.Get = operation
			case "POST":
				pathItem.Post = operation
			case "PUT":
				pathItem.Put = operation
			case "DELETE":
				pathItem.Delete = operation
			}
			config.Paths[path] = pathItem
		}
	}

	return config
}

// 从处理器获取方法
func getMethodFromHandler(handler *HandlerInfo) string {
	return handler.Method
}

// 从处理器获取请求类型
func getRequestType(handler *HandlerInfo) reflect.Type {
	return handler.Request
}

// 生成参数定义
func generateParameters(requestType reflect.Type, handler *HandlerInfo) ([]Parameter, map[string]Definition) {
	parameters := make([]Parameter, 0)
	definitions := make(map[string]Definition)

	if requestType.Kind() == reflect.Struct {
		for i := 0; i < requestType.NumField(); i++ {
			field := requestType.Field(i)
			fieldType := field.Type
			if fieldType == reflect.TypeOf(gin.Context{}) || fieldType == reflect.PointerTo(reflect.TypeOf(gin.Context{})) {
				continue
			}

			// 获取字段标签
			from := field.Tag.Get("from")
			if from == "" {
				from = FromQuery
			}
			mapping := field.Tag.Get("mapping")
			if mapping == "" {
				mapping = field.Name
			}

			// 创建参数
			param := Parameter{
				Name:        mapping,
				In:          from,
				Required:    field.Tag.Get("required") == "true",
				Type:        getSwaggerType(field.Type),
				Description: field.Tag.Get("description"),
			}

			// 如果是复杂类型，添加 schema
			if isComplexType(field.Type) {
				// 生成定义
				var definitionName string
				if fieldType.Kind() == reflect.Struct {
					definition := generateDefinition(fieldType)
					definitionName = fieldType.Name()
					if fieldType.Kind() == reflect.Ptr {
						definitionName = fieldType.Elem().Name()
					}
					if definitionName == "" {
						definitionName = handler.ApiName + "Request"
					}
					definitions[definitionName] = definition
				}

				param.In = "body"
				param.Schema = &Schema{
					Ref: "#/definitions/" + definitionName,
				}
			}

			parameters = append(parameters, param)
		}
	}

	return parameters, definitions
}

// 生成响应定义
func generateResponses(handler *HandlerInfo) (map[string]Response, map[string]Definition) {
	responseType := handler.Response
	if responseType.Kind() == reflect.Ptr {
		responseType = responseType.Elem()
	}

	if !isComplexType(responseType) {
		return map[string]Response{
			"200": {
				Description: "successful operation",
				Schema: Schema{
					Type: getSwaggerType(responseType),
				},
			},
		}, map[string]Definition{}
	}

	definition := generateDefinition(responseType)

	definitionName := responseType.Name()
	if definitionName == "" {
		definitionName = handler.ApiName + "Response"
	}

	return map[string]Response{
			"200": {
				Description: "successful operation",
				Schema: Schema{
					Type: "object",
					Ref:  "#/definitions/" + definitionName,
				},
			},
			"400": {
				Description: "Invalid input",
				Schema: Schema{
					Type: "object",
				},
			},
		}, map[string]Definition{
			definitionName: definition,
		}
}

// 生成定义
func generateDefinition(requestType reflect.Type) Definition {
	definition := Definition{
		Type:       "object",
		Properties: make(map[string]Property),
	}

	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}

	for i := 0; i < requestType.NumField(); i++ {
		field := requestType.Field(i)
		fieldName, property := generateProperty(field)

		// 处理嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			property.Type = "object"
			property.Properties = make(map[string]Property)
			nestedDef := generateDefinition(field.Type)
			property.Properties = nestedDef.Properties
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			property.Type = "object"
			property.Properties = make(map[string]Property)
			nestedDef := generateDefinition(field.Type.Elem())
			property.Properties = nestedDef.Properties
		} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			property.Type = "array"
			property.Items = &Schema{
				Type:       "object",
				Properties: make(map[string]Property),
			}
			nestedDef := generateDefinition(field.Type.Elem())
			property.Items.Properties = nestedDef.Properties
		}

		definition.Properties[fieldName] = property
	}

	return definition
}

func generateProperty(field reflect.StructField) (string, Property) {
	fieldName := field.Name
	if field.Tag.Get("json") != "" {
		fieldName = field.Tag.Get("json")
	}
	property := Property{
		Type:        getSwaggerType(field.Type),
		Description: field.Tag.Get("description"),
	}

	// 处理枚举值
	if enum := field.Tag.Get("enum"); enum != "" {
		property.Enum = strings.Split(enum, ",")
	}

	return fieldName, property
}

// 获取 Swagger 类型
func getSwaggerType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Struct, reflect.Ptr, reflect.Map:
		return "object"
	case reflect.Slice, reflect.Array:
		return "array"
	default:
		panic("unsupported type: " + t.String())
	}
}

func getSummaryFromHandler(handler interface{}) string {
	if h, ok := handler.(IWithSummary); ok {
		return h.GetSummary()
	}
	return ""
}

// 判断是否是复杂类型
func isComplexType(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || t.Kind() == reflect.Slice
}
