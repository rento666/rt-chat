package toolapi

import (
	"rt-chat/internal/api/thirdapi"
	"rt-chat/internal/middleware"
	"rt-chat/pkg/tools/tcaptcha"

	"github.com/gin-gonic/gin"
)

func Api(v *gin.RouterGroup) {
	// 'user包下的api' 浏览器地址：localhost:8080/v1/tool
	tool := v.Group("/tool")
	{
		// '获取(生成)验证码' 浏览器地址：localhost:8080/v1/util/go_captcha_data
		tool.GET("/go_captcha_data", tcaptcha.ApiCreateCaptcha)
		// 校验验证码 浏览器地址：localhost:8080/v1/util/go_captcha_check_data
		tool.POST("/go_captcha_check_data", tcaptcha.ApiCheckCaptcha)
		// 需要token中间件
		tool.GET("/60s", middleware.AuthRequired(), thirdapi.Api60s)
	}
}
