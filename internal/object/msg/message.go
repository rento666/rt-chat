package msg

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"rt-chat/internal/global"
	"rt-chat/internal/object/msg/msgview"
)

// Message 消息
type Message struct {
	gorm.Model
	FormId   uint   `gorm:"comment:发送者id;"`
	TargetId uint   `gorm:"comment:接收者id 用户basic_id 群聊group_id 所有人;"`
	Type     int    `gorm:"comment:消息类型 群聊 私聊 广播;"`
	Media    int    `gorm:"comment:消息类型 文字 图片 音频;"`
	Content  []byte `gorm:"comment:消息内容;"`
}

func (message *Message) TableName() string {
	return "messages"
}

// ConvertMessageVoToMessage 从MessageVo改为message
func ConvertMessageVoToMessage(messageVo msgview.MessageVo) (Message, error) {
	var message Message
	err := json.Unmarshal(messageVo.Content, &message)
	if err != nil {
		log.Fatalln("messageVo to message fail: ", err)
		return message, err
	}
	return message, nil
}

// 增删改查-----------

// AddMsg 添加消息到数据库
func (message *Message) AddMsg() error {

	err := global.DB.Create(message).Error
	if err != nil {
		log.Fatalln("add msg fail: ", err)
		return err
	}
	return nil
}

// 删除消息（撤回？）
// 暂时不写

// 修改消息（撤回+重新编辑？）
// 暂时不写

// FindMsgByFormId 查找消息、需要知道根据什么查：发送人id 排序问题：按时间顺序 返回值：消息列表
func FindMsgByFormId(fid uint) ([]Message, error) {

	var msgs []Message

	err := global.DB.Where("form_id = ?", fid).Order("created_at ASC").Find(&msgs).Error

	if err != nil {
		log.Fatalln("find msg by FormId fail: ", err)
		return nil, err
	}
	return msgs, nil
}

// FindMsgByTargetId 查找消息、需要知道根据什么查：接收者id 排序问题：按时间顺序 返回值：消息列表
func FindMsgByTargetId(tid uint) ([]Message, error) {

	var msgs []Message

	err := global.DB.Where("target_id = ?", tid).Order("created_at ASC").Find(&msgs).Error

	if err != nil {
		log.Fatalln("find msg by TargetId fail: ", err)
		return nil, err
	}
	return msgs, nil
}
