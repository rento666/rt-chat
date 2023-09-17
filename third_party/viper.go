package third_party

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// InitConfig -> viper的初始化
/*
放到init包下，然后就能通过viper.GetString(key)来获取配置文件中的内容了。
*/
func InitConfig() {
	viper.SetConfigName("config") // 读取名为config的配置文件，后缀名为yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // 在当前文件夹下寻找
	err := viper.ReadInConfig()     // 读取配置
	if err != nil {
		fmt.Println("viper初始化失败...")
		log.Fatalln(err)
	} else {
		fmt.Println("viper初始化成功...")
	}

}
