package main

import (
	"fmt"
	"iz2go/api_gen"
)

func main() {
	r := api_gen.GetGin()
	fmt.Println("服务器启动在 :8080 端口...")
	r.Run(":8083")
}
