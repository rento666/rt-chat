package test

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"rt-chat/internal/object/user"
	"testing"
)

func TestG(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3307)/rc?parseTime=True&loc=Local&charset=utf8"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	//db.AutoMigrate(&user.Basic{})
	db.AutoMigrate(&user.Auth{})
	////增
	//db.Create(&user.Basic{Uuid: "1371546845", Nickname: "wrt", Avatar: "touXian",
	//	Status: true, Phone: "iphone 13 pro max 1tb 远峰蓝",
	//	Email: "rento163@163.com", Address: ""})
	//var product user.Basic
	////查
	//db.First(&product, 1)
	////改
	//db.Model(&product).Update("Nickname", "atu")
	////删
	//db.Delete(&product, 1)
}
