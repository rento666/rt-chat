package user

import (
	"gorm.io/gorm"
	"rt-chat/internal/object/client"
)

// Basic 用户基础信息
type Basic struct {
	gorm.Model
	Uuid     string        `gorm:"comment:用户唯一标识;"`
	Nickname string        `gorm:"comment:用户昵称;"`
	Avatar   string        `gorm:"comment:用户头像;"`
	Status   bool          `gorm:"comment:用户状态;"`
	Phone    string        `gorm:"comment:用户手机号;"`
	Email    string        `gorm:"comment:用户邮箱;"`
	Address  string        `gorm:"comment:用户地址;"`
	AuthList []Auth        `gorm:"foreignKey:BasicId"` // 与auth表关联。
	Client   client.Client `gorm:"foreignKey:BasicId"`
}

func (userBasic *Basic) TableName() string {
	return "user_basic"
}
