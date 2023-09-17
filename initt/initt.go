package initt

import (
	"fmt"
	"rt-chat/internal/global"
	"rt-chat/internal/object/user"
	"rt-chat/internal/router"
	"rt-chat/third_party"
)

func Init() {
	third_party.InitConfig()

	global.DB = third_party.InitMysql()
	// db自动迁移模型 如果表不存在则创建表，结构变化则适配变化
	err := global.DB.AutoMigrate(&user.Basic{}, &user.Auth{})
	if err != nil {
		fmt.Println("db.AutoMigrate-error:", err)
		return
	}
	global.ROUTER = router.Router()

	global.REDIS = third_party.InitRedis()

	//tencrypt.GenerateRSAKey(2048)
}
