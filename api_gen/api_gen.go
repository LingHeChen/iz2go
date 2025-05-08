package api_gen

import (
	mod0 "iz2go/routes/api/auth"
	mod1 "iz2go/routes/api/user"

	"github.com/gin-gonic/gin"
)

func GetGin() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Handle(mod0.Method, "/api/auth/Login", (&mod0.Login{}).Execute)
	r.Handle(mod1.Method, "/api/user/GetUserById", (&mod1.GetUserById{}).Execute)
	return r
}
