package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const codeTemplate = `package api_gen

import (
	"github.com/LingHeChen/iz2go/pkg/core"

	{{- range $index, $value := .Routes}}
	mod{{ $index }} "{{ $value.ImportPath }}"
	{{- end}}
)

func InitRoutes() {
	{{- range $index, $value := .Routes}}
	// Register {{.Path}}
	{
		api := &mod{{$index}}.{{.ApiName}}{}
		handlerInfo := iz2go.BuildHandler(api)
		if handlerInfo != nil {
			iz2go.RegisterRoute("{{.Path}}", handlerInfo)
		}
	}
	{{- end}}
}
`

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

func getRoutes(rootPath, modulePath string) []Route {
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

func initPath() (string, string) {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	modulePath, err := parseModulePath(filepath.Join(rootPath, "go.mod"))
	if err != nil {
		log.Fatal("解析 go.mod 失败:", err)
	}
	return rootPath, modulePath
}

func runCmdGen(c *cobra.Command, args []string) {
	rootPath, modulePath := initPath()
	templates := template.Must(template.New("code").Parse(codeTemplate))
	var buf bytes.Buffer
	templates.Execute(&buf, struct {
		Routes []Route
	}{
		Routes: getRoutes(rootPath, modulePath),
	})
	if err := os.MkdirAll(rootPath+"/api_gen", 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(rootPath+"/api_gen/api_gen.go", buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("API 路由已生成")
}

func startServer(filePath string) {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("go", "run", rootPath+"/"+filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runCmdRun(c *cobra.Command, args []string) {
	runCmdGen(c, args)
	filePath := args[0]
	startServer(filePath)
}

var rootCmd = &cobra.Command{
	Use:   "iz2go",
	Short: "iz2go is a tool for generating API routes",
}

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "gen is a tool for generating API routes",
	Run:   runCmdGen,
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "run is a tool for running API routes",
	Run:   runCmdRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cmdGen)
	rootCmd.AddCommand(cmdRun)
}

func main() {
	Execute()
}
