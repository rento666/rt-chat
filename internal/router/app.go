package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"rt-chat/internal/api/system"
	"rt-chat/internal/api/toolapi"
	"rt-chat/internal/object/ws/clientview"
)

func Router() *gin.Engine {
	r := gin.New()
	// 全局中间件
	// 使用 Logger 中间件
	r.Use(gin.Logger())
	// 使用 Recovery 中间件
	r.Use(gin.Recovery())

	// v1版本,项目可多版本共存
	v1 := r.Group("/v1")
	{
		// 将公共api(tools包下的api)挂载到"/v1"上
		toolapi.Api(v1)
		// 将user挂载到"/v1"上,现在v1版本可以访问user包中的api啦
		system.Api(v1)
		// ......

		// 挂载swagger
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// 挂载websocket
		v1.GET("/ws", clientview.WebSocketStart)
	}
	fmt.Println("router初始化成功...")
	return r
}
