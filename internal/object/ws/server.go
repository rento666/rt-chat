package ws

import (
	"log"
	"rt-chat/internal/constants"
	"rt-chat/internal/object/msg/msgview"
	"sync"
	"time"
)

// SktServer 这个不在global中使用，是因为避免循环导包，因为我们还需要将消息存到数据库
var SktServer *Server
var mu sync.Mutex

type Server struct {
	// 在线的客户端列表
	OnlineMap map[uint]*Client
	MapLock   sync.RWMutex
	// 消息广播的chan
	Message chan msgview.Messages
	// 心跳时间，默认为7秒（在NewServer处定义）
	HeartTime time.Duration
}

// InitServer 在init中使用，提供全局server
func InitServer() *Server {
	server := NewServer()
	server.Start()
	return server
}

func NewServer() *Server {
	mu.Lock()
	if SktServer == nil {
		SktServer = &Server{
			OnlineMap: make(map[uint]*Client),
			Message:   make(chan msgview.Messages),
			HeartTime: time.Second * 7,
		}
	}
	mu.Unlock()
	return SktServer
}

// ListenMessage 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (server *Server) ListenMessage() {
	for {
		msgData := <-server.Message
		//将msg发送给全部的在线User
		server.MapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.Msg <- msgData
		}
		server.MapLock.Unlock()
	}
}

// BroadCast 广播消息的方法
func (server *Server) BroadCast(ct *Client, content string) {
	vo := msgview.MessageVo{
		FormId:   ct.UserAuth.BasicId,
		TargetId: constants.ALL,
		Type:     constants.OneToMore,
		Media:    constants.TEXT,
		Content:  []byte(content),
	}
	data, err := msgview.ConvertMessageVoToByte(vo)
	if err != nil {
		log.Println("messageVo to []byte fail: ", err)
		return
	}
	sendMsg := msgview.Messages{
		MsgType: 1,
		Data:    data,
	}
	server.Message <- sendMsg
}

func (server *Server) Handler(ct *Client) {

	// 客户端上线
	ct.Online()

	isLive := make(chan bool)
	defer close(isLive)

	// 另开一个协程，负责读取发送消息，然后保持活跃状态
	go func() {
		for {
			mt, data, err := ct.Conn.ReadMessage()
			// 连接断开的情况下，这里会出错
			if err != nil {
				log.Println("read fail: ", err)
				return
			}
			msgs := &msgview.Messages{
				MsgType: mt,
				Data:    data,
			}
			// 客户端去处理读取的消息
			ct.DoMessage(msgs)
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		case <-time.After(server.HeartTime):
			close(ct.Msg)
			return
		}
	}
}

func (server *Server) Start() {
	// 监听广播消息
	go server.ListenMessage()
}
