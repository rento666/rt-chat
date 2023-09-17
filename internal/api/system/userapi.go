package system

import (
	"github.com/gin-gonic/gin"
	"rt-chat/internal/middleware"
	"rt-chat/internal/object/user/userview"
)

func Api(v *gin.RouterGroup) {
	// 'user包下的api' 浏览器地址：localhost:8002/v1/user
	user := v.Group("/user")
	{
		user.POST("/register", userview.Register)
		user.GET("/register_code", userview.RegisterCode)
		user.GET("/judge", userview.JudgeUserExist)
		user.POST("/login", userview.Login)
		user.GET("/list", middleware.AuthRequired(), userview.List)
	}
}
