package constants

// object/msg/message
const (
	NOBODY = iota - 1 // 接收者：没有人，此种情况为保持心跳的Ping，无实际作用
	ALL               // 所有在线的用户

	OneToOne  // 私聊
	OneToMany // 群聊
	OneToMore // 广播

	TEXT    = iota - 4 // 文字
	PICTURE            // 图片
	AUDIO              // 语音
	VIDEO              // 视频

)
