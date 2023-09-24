package client

import (
	"gorm.io/gorm"
	"time"
)

type Client struct {
	gorm.Model
	BasicId       uint      `gorm:"comment:basic表的id"`
	ClientIp      string    `gorm:"客户端Ip"`
	ClientPort    string    `gorm:"客户端端口"`
	LoginTime     time.Time `gorm:"登录时间"`
	HeartbeatTime time.Time `gorm:"心跳时间"`
	LoginOutTime  time.Time `gorm:"下线时间"`
	IsLogout      bool      `gorm:"是否下线"`
	DeviceInfo    string    `gorm:"设备信息"`
}

func (client *Client) TableName() string {
	return "sys_client"
}
