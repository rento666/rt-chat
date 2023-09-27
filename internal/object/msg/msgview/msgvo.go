package msgview

import (
	"encoding/json"
	"log"
)

// MessageVo 消息
type MessageVo struct {
	FormId   uint   `comment:"发送者id;"`
	TargetId uint   `comment:"接收者id 用户basic_id 群聊group_id;"`
	Type     int    `comment:"消息类型 群聊 私聊 广播;"`
	Media    int    `comment:"消息类型 文字 图片 音频;"`
	Content  []byte `comment:"消息内容;"`
}

type Messages struct {
	MsgType int
	Data    []byte
}

// ConvertMessagesToMessageVo 从data里解包messageVo
func ConvertMessagesToMessageVo(data []byte) (MessageVo, error) {
	var messageVo MessageVo
	err := json.Unmarshal(data, &messageVo)
	if err != nil {
		log.Fatalln("data to messageVo fail: ", err)
		return messageVo, err
	}
	return messageVo, nil
}

// ConvertMessageVoToByte messageVo 转为 []byte的data
func ConvertMessageVoToByte(messageVo MessageVo) ([]byte, error) {
	data, err := json.Marshal(messageVo)
	if err != nil {
		log.Fatalf("Marshal failed: %v", err)
	}
	return data, nil
}
