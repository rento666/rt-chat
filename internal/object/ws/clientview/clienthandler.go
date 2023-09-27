package clientview

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"rt-chat/internal/middleware"
	"rt-chat/internal/object/ws"
)

// 防止跨域伪造请求
var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketStart(c *gin.Context) {
	/*
		前端发送http请求，升级为websocket：
			1.建立链接
			2.创建一个客户端对象，保持这个链接的信息
			3.通过全局的服务端进行处理业务
	*/
	con, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer con.Close()

	auth, err := middleware.GetAuthByToken(c)
	if err != nil {
		log.Println("token fail")
		return
	}
	// 创建一个客户端对象
	ct := ws.NewClient(c.ClientIP(), auth, con)

	// websocket 处理该客户端的业务
	ws.SktServer.Handler(ct)
}
