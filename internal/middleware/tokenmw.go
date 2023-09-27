package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rt-chat/internal/object/user"
	"rt-chat/internal/object/user/userview"
	"rt-chat/internal/result"
	"rt-chat/pkg/tools/tjwt"
	"strconv"
	"strings"
	"time"
)

// AuthRequired 路由请求中间件，前端必须把token放在请求头上，对服务器进行请求验证token成功后，才能访问后续的请求路由
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 authorization header：获取前端传过来的信息的
		tokenString := c.GetHeader("Authorization")

		//验证前端传过来的token格式，不为空，开头为Bearer
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(401, result.Json(http.StatusUnauthorized, "权限不足", nil))
			c.Abort()
			return
		}
		//验证通过，提取有效部分（除去Bearer)
		tokenString = tokenString[7:] //截取字符
		//解析token
		claims, err := tjwt.CheckToken(tokenString)
		//解析失败||解析后的token无效
		if err != nil {
			c.JSON(401, result.Json(http.StatusUnauthorized, fmt.Sprintf("%v", err), nil))
			c.Abort()
			return
		}
		if !claims.VerifyExpiresAt(time.Now(), false) {
			c.JSON(406, result.Json(http.StatusNotAcceptable, "token已过期", nil))
			c.Abort()
			return
		}

		token, err := tjwt.RefreshToken(claims)
		if err != nil {
			c.JSON(401, result.Json(http.StatusUnauthorized, fmt.Sprintf("token刷新失败:%v", err), nil))
			c.Abort()
			return
		}
		// 刷新token
		c.Header("NEW-TOKEN", token)

		uid := claims.User
		userid, err := strconv.ParseUint(uid, 10, 32)
		if err != nil {
			fmt.Println("token的uid转换失败:", err)
			return
		}
		b, u := userview.JudgeExistFunc(user.Auth{BasicId: uint(userid)})
		if !b {
			c.JSON(401, result.Json(http.StatusUnauthorized, "权限不足", nil))
			c.Abort()
			return
		}
		// 可以在跨中间件获取当前用户auth信息。账密、UUID等
		c.Set("user_auth", u)
		c.Next()
	}
}

func GetAuthByToken(c *gin.Context) (user.Auth, error) {
	// 从token那里获取user_auth
	value, exists := c.Get("user_auth")
	if !exists {
		fmt.Println("token isn't exists!")
		return user.Auth{}, errors.New("无法从token处获取用户信息")
	}
	auth := value.(user.Auth)
	return auth, nil
}
