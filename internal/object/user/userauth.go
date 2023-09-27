package user

import (
	"gorm.io/gorm"
)

// IdentityType类型
const (
	EMAIL  = "email"
	PHONE  = "phone"
	QQ     = "qq"
	WECHAT = "wechat"
	NAME   = "username"
)

type Auth struct {
	gorm.Model
	BasicId      uint   `gorm:"comment:basic表的id"`
	IdentityType string ` gorm:"type:varchar(20); comment:用户标识类型(邮箱、手机号、token);"`
	Identifier   string ` gorm:"comment:用户登录标识;"`
	Credential   string ` gorm:"comment:用户密钥;"`
}

func (userAuth *Auth) TableName() string {
	return "user_auth"
}
