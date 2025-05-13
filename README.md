# iz2go

一个用于路由生成的工具

## 使用方法

### 1. 安装依赖包和生成工具

```shell
go get "github.com/LingHeChen/iz2go"
go install "github.com/LingHeChen/iz2go/cmd/iz2go"
```

### 2. 编写路由文件

在项目根目录下创建`/routes`目录，并根据需要添加目录结构(即路由结构)

```plainText
routes
├── auth
│   └── Login.go // /auth/Login
└── path_test
    └── Handler@*path.go // /auth/*path
```

使用@可以重写路径

编写路由处理器文件，以`/auth/Login`接口为例

```golang
package auth

import (
	"iz2go_test/middlewires"
	"iz2go_test/services"

	iz2go "github.com/LingHeChen/iz2go/pkg/core"
	"github.com/gin-gonic/gin"
)

type Login struct {
	*iz2go.Post
	LoginService *services.LoginService
}

func (api *Login) Init() {
	api.LoginService = &services.LoginService{}
}

func (api *Login) Execute(request struct {
	Ctx  *gin.Context
	Body struct {
		Username string `json:"username" description:"用户名"`
		Password string `json:"password" description:"密码"`
	}
}) (struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}, iz2go.IError) {
	token, err := api.LoginService.Login(request.Body.Username, request.Body.Password)
	if err != nil {
		return struct {
			Message string `json:"message"`
			Token   string `json:"token"`
		}{
			Message: "Unauthorized",
			Token:   token,
		}, nil
	}
	return struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "Unauthorized",
		Token:   token,
	}, nil
}

func (api *Login) Decorators() []iz2go.Decorator {
	return []iz2go.Decorator{
		middlewires.RequireRoles([]string{"admin"}),
	}
}
```

### 3. 生成路由代码并运行

在项目的根目录下目录下运行以下命令

```shell
iz2go gen
```

在入口文件(如`main.go`)中添加如下代码：

```golang
package main

import (
	"iz2go_test/api_gen"
	"net/http"

	iz2go "github.com/LingHeChen/iz2go/pkg/core"
	"github.com/gin-gonic/gin"
)

func main() {
	api_gen.InitRoutes()
	iz2go.RegisterErrorHook(func(ctx *gin.Context, err iz2go.IError) (iz2go.IError, bool) {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": err.GetCode()})
		return nil, true
	})
	r := iz2go.Default()
	r.RenderSwagger(&iz2go.SwaggerRenderConfig{
		Info: &iz2go.Info{
			Title:       "API Documentation",
			Description: "API Documentation for the project",
			Version:     "1.0.0",
		},
	})
	r.Run(":8083")
}
```

然后运行项目即可
或者在添加main代码后，直接运行`iz2go run <入口文件名>`

## 未来计划

* [X]  添加参数的自动绑定
* [X]  添加全局错误处理🪝
* [X]  集成swagger
* [ ]  更完整的swagger支持
* [ ]  添加更好的websocket支持
* [ ]  添加配置文件
