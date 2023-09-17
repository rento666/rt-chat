package result

import "github.com/gin-gonic/gin"

func Json(code int, msg string, data any) map[string]any {
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}
