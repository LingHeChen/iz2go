package iz2go

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

var errorHooks []func(ctx *gin.Context, err IError) (IError, bool)
var successHooks []func(ctx *gin.Context, response gin.H) (gin.H, bool)

func OnError(c *gin.Context, err IError) {
	abort := false
	hooks := errorHooks
	slices.Reverse(hooks)
	for _, hook := range hooks {
		err, abort = hook(c, err)
		if abort {
			return
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": err.GetCode(), "message": err.GetMessage()})
}

func OnSuccess(c *gin.Context, response gin.H) {
	abort := false
	hooks := successHooks
	slices.Reverse(hooks)
	for _, hook := range hooks {
		response, abort = hook(c, response)
		if abort {
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func RegisterErrorHook(hook func(ctx *gin.Context, err IError) (IError, bool)) {
	errorHooks = append(errorHooks, hook)
}

func RegisterSuccessHook(hook func(ctx *gin.Context, response gin.H) (gin.H, bool)) {
	successHooks = append(successHooks, hook)
}
