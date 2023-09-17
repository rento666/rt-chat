package tcaptcha

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rt-chat/internal/result"

	_ "rt-chat/docs"
)

// ApiCreateCaptcha 生成验证码的api
// @Tags 生成验证码
// @Success 200 {object} Ct "captcha"
// @Router /v1/tool/go_captcha_data [get]
func ApiCreateCaptcha(c *gin.Context) {
	ct := CreateCode()
	code := http.StatusOK
	msg := "获取人机验证数据成功"
	if ct == nil {
		code = http.StatusForbidden
		msg = "获取人机验证数据失败"
	}
	c.JSON(http.StatusOK, result.Json(code, msg, ct))
}

// ApiCheckCaptcha 校验验证码的api
// @Tags 校验验证码
// @Accept json
// @Produce json
// @Param data body Che true "{dots:"",key:""}"
// @Success 200 {bool} isValid
// @Router /v1/tool/go_captcha_check_data [post]
func ApiCheckCaptcha(c *gin.Context) {
	// 获取前端传来的dots、key，默认为字符串nil

	dots := c.DefaultPostForm("dots", "nil")
	key := c.DefaultPostForm("key", "nil")
	msg, isTrue, code := CheckCode(dots, key)
	// 返回结构体
	c.JSON(http.StatusOK, result.Json(code, msg, isTrue))
}
