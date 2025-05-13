package iz2go

import (
	"net/http"
	"reflect"
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

func OnSuccess(c *gin.Context, response interface{}) {
	abort := false
	hooks := successHooks
	slices.Reverse(hooks)
	for _, hook := range hooks {
		if reflect.TypeOf(response) == reflect.TypeOf(reflect.TypeOf(hook).In(1)) {
			ret := reflect.ValueOf(hook).Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(response)})
			response = ret[0].Interface()
			if !ret[1].IsNil() {
				abort = ret[1].Interface().(bool)
			}
			if abort {
				return
			}
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
