package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Route struct {
	ImportPath string
	Path       string
	ApiName    string
}

// 从 go.mod 解析模块路径
func parseModulePath(goModPath string) (string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(line[len("module "):]), nil
		}
	}

	return "", fmt.Errorf("未找到 module 声明")
}

func getRoutes() []Route {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	modulePath, err := parseModulePath(filepath.Join(rootPath, "go.mod"))
	if err != nil {
		log.Fatal("解析 go.mod 失败:", err)
	}

	routes := []Route{}
	filepath.Walk(rootPath+"/routes", func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		importPath := strings.Replace(filePath, rootPath, "", 1)
		importPath = strings.Replace(importPath, ".go", "", 1)
		apiName := strings.Split(importPath, "/")[len(strings.Split(importPath, "/"))-1]
		importPath = strings.Join(strings.Split(importPath, "/")[:len(strings.Split(importPath, "/"))-1], "/")

		path := strings.Replace(importPath, "/routes", "", 1)

		routes = append(routes, Route{
			ImportPath: modulePath + importPath,
			Path:       path + "/" + apiName,
			ApiName:    apiName,
		})

		return nil
	})
	return routes
}

func main() {
	const codeTemplate = `package api_gen

import (
	"github.com/gin-gonic/gin"
	{{- range $index, $value := .Routes}}
	mod{{ $index }} "{{ $value.ImportPath }}"
	{{- end}}
)

func GetGin() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	{{- range $index, $value := .Routes}}
	// Register {{.Path}}
	{
		api := &mod{{$index}}.{{.ApiName}}{}
		r.Handle(mod{{$index}}.Method, "{{.Path}}", api.Execute)
	}
	{{- end}}
	return r
}
`
	templates := template.Must(template.New("code").Parse(codeTemplate))
	var buf bytes.Buffer
	templates.Execute(&buf, struct {
		Routes []Route
	}{
		Routes: getRoutes(),
	})
	os.WriteFile("api_gen/api_gen.go", buf.Bytes(), 0644)

	fmt.Println("API 路由已生成")
}
