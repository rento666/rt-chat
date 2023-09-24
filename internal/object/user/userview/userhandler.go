package userview

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rt-chat/internal/object/user"
	"rt-chat/internal/result"
	"rt-chat/pkg/tools/tencrypt"
)

// 用户的业务处理

// Register 用户注册api方法
// @Summary 新增用户
// @Tags 用户模块
// @Accept json
// Produce json
// @Param data body RegisterVo true "{IdentityType,Identifier,Credential,Code}"
// @Success 200 {string} token
// @Router /v1/user/register [post]
func Register(c *gin.Context) {
	var form RegisterVo
	code := http.StatusUnsupportedMediaType
	msg := "注册失败，填写数据有误"
	var data any

	if c.ShouldBind(&form) == nil && form.IdentityType != "" &&
		form.Identifier != "" && form.Credential != "" && form.Code != "" {
		token, err := RegisterFunc(form)
		if err != nil {
			code = http.StatusNoContent
			msg = fmt.Sprintf("%v", err)
		} else {
			code = http.StatusOK
			msg = "注册成功"
			data = token
		}
	}
	c.JSON(http.StatusOK, result.Json(code, msg, data))
}

// JudgeUserExist 判断用户是否存在api(当使用注册、登录(?)功能时使用)  返回rsa的公钥
// @Summary 判断用户是否存在
// @Tags 用户模块
// @Param Identifier query string true "账号"
// @Param IdentityType query string true "账号类型"
// @Success 200 {bool} isValid
// @Router /v1/user/judge [get]
func JudgeUserExist(c *gin.Context) {
	code := http.StatusUnsupportedMediaType
	msg := "请求参数有误"
	var data any
	// 如果为true，则说明传入了手机号\邮箱\assess_token
	identifier := c.DefaultQuery("Identifier", "nil")
	identifierType := c.DefaultQuery("IdentityType", "nil")
	if identifier != "nil" && identifierType != "nil" {
		e, _ := JudgeExistFunc(user.Auth{IdentityType: identifierType, Identifier: identifier})
		if !e {
			code = http.StatusNoContent
			// 不存在，可以进行注册，不能登录
			msg = "用户不存在"
		} else {
			code = http.StatusOK
			// 存在，可以登录，不能再次注册
			msg = "用户已存在"
		}
		// 获取公钥
		publicKey, err := tencrypt.GetPublicKey()
		data = publicKey
		if err != nil {
			fmt.Println("获取rsa公钥失败:", err)
			data = ""
		}
	}
	c.JSON(http.StatusOK, result.Json(code, msg, data))
}

// Login 用户登录，先进行rsa解密，得到账号类型、账号、密码；再进行aes解密，得到账号
// @Summary 用户登录
// @Tags 用户模块
// @Accept json
// Produce json
// @Param data body LoginVo true "{Info}"
// @Success 200 {string} token
// @Router /v1/user/login [post]
func Login(c *gin.Context) {
	var form LoginVo
	code := http.StatusUnauthorized
	msg := "请求参数有误"
	var data any
	if c.ShouldBind(&form) == nil && form.Info != "" {
		token, err := LoginFunc(form)
		if err != nil {
			code = http.StatusNoContent
			msg = fmt.Sprintf("%v", err)
		} else {
			code = http.StatusOK
			msg = "登录成功"
			data = token
		}
	}
	c.JSON(http.StatusOK, result.Json(code, msg, data))
}

// RegisterCode 注册验证码（邮箱、手机号）
// @Summary 注册验证码
// @Tags 用户模块
// Produce json
// @Param Identifier query string true "账号"
// @Param IdentityType query string true "账号类型"
// @Success 200 {string} msg
// @Router /v1/user/register_code [get]
func RegisterCode(c *gin.Context) {
	code := http.StatusBadRequest
	msg := "获取验证码失败，填写数据有误"
	var data any

	identifier := c.DefaultQuery("Identifier", "nil")
	identityType := c.DefaultQuery("IdentityType", "nil")

	if identifier != "nil" && identityType != "nil" {
		if !IdentifierCodeFunc(IdentifierCode{IdentityType: identityType, Identifier: identifier}) {
			msg = "不支持的账号类型，发送验证码失败！"
		} else {
			code = http.StatusOK
			msg = "发送成功，验证码30分钟有效"
		}
	}
	c.JSON(http.StatusAccepted, result.Json(code, msg, data))
}

// List 获取用户列表，需要有token才能用，所以必须要有中间件判断
// @Summary 用户列表
// @Tags 用户模块
// Produce json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param pageNum query int false "页码数量"
// @Param pageSize query int false "分页大小"
// @Param keyword query string false "模糊搜索关键字"
// @Param desc query bool false "是否反向搜索"
// @Success 200 {object} UserList "list"
// @Router /v1/user/list [get]
func List(c *gin.Context) {
	var p result.Page
	code := http.StatusUnauthorized
	msg := "请求参数有误"
	var data any
	if c.ShouldBindQuery(&p) == nil {
		// {pageNum int 页码数量,pageSize int 分页大小,keyword string 模糊搜索关键字,desc bool 是否反向搜索}
		if p.PageNum <= 0 {
			p.PageNum = 1
		}
		list, err := GetUserListFunc(p)
		if err != nil {
			msg = fmt.Sprintf("%v", err)
		} else {
			code = http.StatusOK
			msg = "获取成功"
			data = list
		}
	}
	c.JSON(http.StatusOK, result.Json(code, msg, data))
}
