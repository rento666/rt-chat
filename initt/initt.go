package initt

import (
	"fmt"
	"rt-chat/internal/global"
	"rt-chat/internal/object/group"
	"rt-chat/internal/object/msg"
	"rt-chat/internal/object/user"
	"rt-chat/internal/object/ws"
	"rt-chat/internal/router"
	"rt-chat/third_party"
)

func Init() {
	third_party.InitConfig()

	global.DB = third_party.InitMysql()
	// db自动迁移模型 如果表不存在则创建表，结构变化则适配变化 新增数据库表时需要修改
	err := global.DB.AutoMigrate(&user.Basic{}, &user.Auth{},
		&msg.Message{}, &group.Group{})
	if err != nil {
		fmt.Println("db.AutoMigrate-error:", err)
		return
	}
	global.ROUTER = router.Router()

	global.REDIS = third_party.InitRedis()

	ws.SktServer = ws.InitServer()

	//tencrypt.GenerateRSAKey(2048)
}
