package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"rt-chat/internal/constants"
	"rt-chat/internal/object/msg"
	"rt-chat/internal/object/msg/msgview"
	"rt-chat/internal/object/user"
)

type Client struct {
	Ip       string                // 客户端IP
	Msg      chan msgview.Messages // 消息
	UserAuth user.Auth             // 由token获取的user_auth
	Conn     *websocket.Conn       // 客户端的websocket连接
	Server   *Server
}

func NewClient(ip string, ua user.Auth, conn *websocket.Conn) *Client {
	return &Client{
		Ip:       ip,
		Msg:      make(chan msgview.Messages),
		UserAuth: ua,
		Conn:     conn,
	}
}

// Online 客户端的上线业务
func (ct *Client) Online() {

	//用户上线,将用户加入到onlineMap中
	ct.Server.MapLock.Lock()
	ct.Server.OnlineMap[ct.UserAuth.BasicId] = ct
	ct.Server.MapLock.Unlock()

	//广播当前用户上线消息
	ct.Server.BroadCast(ct, "上线了")
}

// Offline 客户端的下线业务
func (ct *Client) Offline() {
	//用户下线,将用户从onlineMap中删除
	ct.Server.MapLock.Lock()
	delete(ct.Server.OnlineMap, ct.UserAuth.BasicId)
	ct.Server.MapLock.Unlock()

	//广播当前用户上线消息
	ct.Server.BroadCast(ct, "下线了")
}

// SendMsg 给当前客户端发送消息
func (ct *Client) SendMsg(msgs *msgview.Messages) {
	err := ct.Conn.WriteMessage(msgs.MsgType, msgs.Data)
	if err != nil {
		log.Println("send msg to client [ip:"+ct.Ip+"] fail: ", err)
		return
	}
}

// DoMessage 客户端处理消息的业务
func (ct *Client) DoMessage(msgs *msgview.Messages) {
	// 形参为读取到的消息结构体
	data := msgs.Data
	// 这个data是前端发来的消息，是JSON类型的、与messageVo结构体结构一样
	vo, err := msgview.ConvertMessagesToMessageVo(data)
	if err != nil {
		log.Println("[client DoMessage] messages to msgVo fail: ", err)
		return
	}
	// messageVo里包含接收者信息，可以进行消息发送！
	switch vo.Type {
	case constants.OneToOne:
		// 私聊，向目标用户（单个）发送消息 从当前客户端的server端查看在线列表（仅限单体服务器上）
		remoteCt, ok := ct.Server.OnlineMap[vo.TargetId]
		if !ok {
			vo.Content = []byte("该用户目前不在线~")
			bts, err := msgview.ConvertMessageVoToByte(vo)
			if err != nil {
				log.Println("[client DoMessage] msgVo to data fail: ", err)
				return
			}
			ct.SendMsg(&msgview.Messages{MsgType: msgs.MsgType, Data: bts})
			return
		}
		remoteCt.SendMsg(msgs)
		log.Printf("success send msg %v to ip:%s from ip:%s \n", string(vo.Content), remoteCt.Ip, ct.Ip)
	case constants.OneToMany:
		// 群聊消息，这个需要从群里依次获取用户id，然后保存到数据库，在线的会接收到推送
	case constants.OneToMore:
		// 广播消息
		ct.Server.BroadCast(ct, string(vo.Content))
	}

	// messageVo转为message，这两者的差别就是：Vo没有gorm.Model这个内嵌结构体
	message, err := msg.ConvertMessageVoToMessage(vo)
	if err != nil {
		log.Println("[client DoMessage] msgVo to message fail: ", err)
		return
	}
	err = message.AddMsg()
	if err != nil {
		log.Println("[client DoMessage] message add to sql fail: ", err)
		return
	}
}
