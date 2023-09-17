package third_party

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func dsn() (dsn string) {
	var path = viper.GetString("mysql.path")
	var port = viper.GetString("mysql.port")
	var config = viper.GetString("mysql.config")
	var username = viper.GetString("mysql.username")
	var password = viper.GetString("mysql.password")
	var dbname = viper.GetString("mysql.dbname")
	return username + ":" + password + "@tcp(" + path + ":" + port + ")/" + dbname + "?" + config
}

// InitMysql -> mysql的初始化
/*
放到init包下的，然后就能通过db.Create()来新增数据（增删改查）了。
*/
func InitMysql() *gorm.DB {
	dsn := dsn()
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	// 自定义日志模板，打印sql语句

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		})

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{Logger: newLogger}); err != nil {
		fmt.Println("mysql初始化失败...")
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetConnMaxLifetime(time.Hour)
		fmt.Println("mysql初始化成功...")
		return db
	}
}

// InitRedis 连接redis
/*
@return rc *redis.Client
*/
func InitRedis() *redis.Client {
	var addr = viper.GetString("redis.addr")
	var password = viper.GetString("redis.password")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	fmt.Println("redis初始化成功...")
	return redisClient
}
