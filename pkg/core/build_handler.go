package iz2go

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	FromQuery  = "query"
	FromPath   = "path"
	FromCtx    = "ctx"
	FromHeader = "header"
)

func BuildHandler(handler interface{}) *HandlerInfo {
	if handler == nil {
		return nil
	}

	executableMethod := CheckExecutable(handler)

	HandleInit(handler)
	method := ParseMethod(handler)
	handlerFunc := ParseHandler(handler, executableMethod)

	return &HandlerInfo{
		Method:   method,
		Handler:  handlerFunc,
		Request:  executableMethod.Type.In(1),
		Response: executableMethod.Type.Out(0),
	}
}

func CheckExecutable(handler interface{}) reflect.Method {
	// 获取handler的类型
	handlerType := reflect.TypeOf(handler)

	// 查找Execute方法
	method, ok := handlerType.MethodByName("Execute")
	if !ok {
		panic("handler must have Execute method")
	}

	// 检查Execute方法的签名
	if method.Type.NumIn() != 2 {
		panic("Execute method must have two and only two parameters")
	}

	// 检查Execute方法的返回值
	if method.Type.NumOut() != 2 {
		panic("Execute method must return two and only two values")
	}
	errorType := method.Type.Out(1)
	if !errorType.Implements(reflect.TypeOf((*IError)(nil)).Elem()) {
		panic("Execute method must return (Any, IError)")
	}
	return method
}

func ParseMethod(handler interface{}) string {
	if h, ok := handler.(IWithMethod); ok {
		return h.GetMethod()
	}
	return "GET"
}

func HandleInit(handler interface{}) {
	if h, ok := handler.(IWithInit); ok {
		h.Init()
	}
}

func CheckDecorator(handler interface{}) ([]Decorator, bool) {
	// 获取handler的类型
	handlerType := reflect.TypeOf(handler)

	// 查找Execute方法
	method, ok := handlerType.MethodByName("Decorators")
	if !ok {
		return nil, false
	}

	// 检查Execute方法的签名
	if method.Type.NumIn() > 0 {
		return nil, false
	}

	// 检查Execute方法的返回值
	if method.Type.NumOut() != 1 ||
		method.Type.Out(0).String() != "[]Decorator" {
		return nil, false
	}
	// 创建值并调用方法获取结果
	handlerValue := reflect.ValueOf(handler)
	decoratorsValue := handlerValue.MethodByName("Decorators").Call(nil)[0]
	return decoratorsValue.Interface().([]Decorator), true
}

func ParseHandler(handler interface{}, executableMethod reflect.Method) gin.HandlerFunc {
	decorators, ok := CheckDecorator(handler)
	handlerFunc := executableMethod.Func
	wrapperedHandlerFunc := WrapperHandlerFunc(reflect.ValueOf(handler), handlerFunc)
	if ok {
		slices.Reverse(decorators)
		for _, decorator := range decorators {
			wrapperedHandlerFunc = decorator(wrapperedHandlerFunc)
		}
		return wrapperedHandlerFunc
	}
	return wrapperedHandlerFunc
}

func WrapperHandlerFunc(handler reflect.Value, handlerFunc reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := ParseRequest(c, handlerFunc.Type().In(1))
		ret := handlerFunc.Call([]reflect.Value{handler, request})
		var err IError
		var ok bool
		response := ret[0].Interface()
		if !ret[1].IsNil() {
			err, ok = ret[1].Interface().(IError)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error must be IError"})
				return
			}
		} else {
			err = nil
		}

		if err != nil {
			OnError(c, err)
			return
		}
		OnSuccess(c, response)
	}
}

func ParseRequest(c *gin.Context, requestType reflect.Type) reflect.Value {
	// 如果是 *gin.Context 类型，直接返回 context
	if requestType == reflect.TypeOf(c) {
		return reflect.ValueOf(c)
	}

	// 创建请求类型的新实例
	request := reflect.New(requestType).Elem()

	// 如果是结构体类型，尝试先从 JSON body 绑定
	if requestType.Kind() == reflect.Struct {

		// 如果绑定失败，继续使用字段解析方式
		for i := 0; i < requestType.NumField(); i++ {
			field := request.Field(i)
			fieldType := requestType.Field(i)

			// 跳过不可设置的字段
			if !field.CanSet() {
				continue
			}

			if fieldType.Type == reflect.TypeOf(c) {
				field.Set(reflect.ValueOf(c))
				continue
			}

			// 获取字段标签
			from := fieldType.Tag.Get("from")
			if from == "" {
				from = FromQuery
			}
			mapping := fieldType.Tag.Get("mapping")
			if mapping == "" {
				mapping = fieldType.Name
			}

			// 根据字段类型设置值
			switch field.Kind() {
			case reflect.String:
				value := getValueFromContext(c, from, mapping)
				field.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if value := getValueFromContext(c, from, mapping); value != "" {
					if val, err := strconv.ParseInt(value, 10, 64); err == nil {
						field.SetInt(val)
					}
				}
			case reflect.Float32, reflect.Float64:
				if value := getValueFromContext(c, from, mapping); value != "" {
					if val, err := strconv.ParseFloat(value, 64); err == nil {
						field.SetFloat(val)
					}
				}
			case reflect.Bool:
				if value := getValueFromContext(c, from, mapping); value != "" {
					if val, err := strconv.ParseBool(value); err == nil {
						field.SetBool(val)
					}
				}
			case reflect.Struct:
				// 处理嵌套结构体 - 这里可以递归调用 ParseRequest
				jsonObj := reflect.New(fieldType.Type).Interface()
				if err := c.ShouldBindJSON(jsonObj); err == nil {
					field.Set(reflect.ValueOf(jsonObj).Elem())
				}
			}
		}
	}

	return request
}

// 辅助函数：根据来源获取值
func getValueFromContext(c *gin.Context, from string, mapping string) string {
	switch from {
	case FromQuery:
		return c.Query(mapping)
	case FromPath:
		return c.Param(mapping)
	case FromCtx:
		if value, exists := c.Get(mapping); exists {
			return fmt.Sprintf("%v", value)
		}
	case FromHeader:
		return c.GetHeader(mapping)
	}
	return ""
}

func buildHandler(handler interface{}) {
	// 获取handler的类型
	handlerType := reflect.TypeOf(handler)

	// 查找Execute方法
	method, ok := handlerType.MethodByName("Execute")
	if !ok {
		panic("handler must have Execute method")
	}

	// 检查Execute方法的签名
	if method.Type.NumIn() < 2 {
		panic("Execute method must have at least one parameter")
	}

	// 检查Execute方法的返回值
	if method.Type.NumOut() != 2 ||
		method.Type.Out(0).String() != "gin.H" ||
		method.Type.Out(1).String() != "error" {
		panic("Execute method must return (gin.H, error)")
	}

	// 获取Execute方法的参数类型
	paramType := method.Type.In(1)

	fmt.Printf("Handler implements Execute with parameter type: %s\n", paramType.String())

	// 这里可以进行其他处理，如存储参数类型信息或动态生成适配器
}
