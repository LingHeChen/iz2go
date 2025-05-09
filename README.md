# iz2go

一个用于路由生成的工具

## 使用方法
### 1. 安装依赖包和生成工具
``` shell
go get "github.com/LingHeChen/iz2go"
go isntall "github.com/LingHeChen/iz2go/cmd/iz2go"
```
### 2. 编写路由文件
在项目根目录下创建`/routes`目录，并根据需要添加目录结构(即路由结构)
``` plainText
routes
├── auth
│   └── Login.go // /auth/Login
└── path_test
    └── Handler@*path.go // /auth/*path
```
使用@可以重写路径

编写路由处理器文件，以`/auth/Login`接口为例
``` golang
package auth

import (
	"iz2go_test/middlewires"
	"iz2go_test/services"

	iz2go "github.com/LingHeChen/iz2go/pkg/core"
	"github.com/gin-gonic/gin"
)

type Login struct {
	*iz2go.Get // 接口方法 (可选，默认为GET)
	LoginService *services.LoginService
}

// 接口初始化 (可选)
func (api *Login) Init() {
	api.LoginService = &services.LoginService{}
}

// 接口入口
func (api *Login) Execute(c *gin.Context) {
	token, err := api.LoginService.Login("admin", "admin")
	if err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized"})
	}
	c.JSON(200, gin.H{"message": "Hello, World!", "token": token})
}

// 装饰器 (可选)
func (api *Login) Decorators() []iz2go.Decorator {
	return []iz2go.Decorator{
		middlewires.RequireRoles([]string{"admin"}),
	}
}
```
### 3. 生成路由代码并运行
在项目的根目录下目录下运行以下命令
``` shell
iz2go gen
```
在入口文件(如`main.go`)中添加如下代码：
``` golang
import (
    "github.com/LingHeChen/iz2go"
    "<模块名>/api_gen"
)

func main() {
    // ...其他代码
    api_gen.InitRoutes()
    r := iz2go.NewRouter() // 会返回一个 *gin.Engine
    r.Run(":<port>")
    // ...其他代码
}
```
然后运行项目即可
或者在添加main代码后，直接运行`iz2go run <入口文件名>`
