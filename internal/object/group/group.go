package group

import (
	"gorm.io/gorm"
	"rt-chat/internal/object/user"
)

type Group struct {
	gorm.Model
	GroupId  uint         `gorm:"comment:群聊ID"`
	Name     string       `gorm:"comment:群聊名称"`
	UserList []user.Basic `gorm:"foreignKey:GroupId"` // 与user_basic表关联
}

func (group *Group) TableName() string {
	return "group"
}
