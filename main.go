package main

import (
	"rt-chat/initt"
	"rt-chat/internal/global"
)

func init() {

	initt.Init()
}

func main() {
	// 启动服务
	global.ROUTER.Run(":8002")
}
