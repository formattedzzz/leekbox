package stream

import (
	"encoding/json"
	"fmt"
	"leekbox/api/auth"
	"leekbox/model"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 设置该客户端下一个消息最大写入时长
	writeWait = 10 * time.Second
	// 允许客户端发出下一个消息的最大时长
	pongWait = 60 * time.Second
	// 客户端发出的心跳检测周期
	// 服务端client对象每ping成功一次会触发PongHandler 延长pongWait 故该值必须小于pongWait
	pingPeriod = (pongWait * 9) / 10
	// 客户端消息的允许最大字节数
	maxMessageSize = 2048
)

const (
	// 接受类型
	MESSAGE_TYPE_LOGIN  = "LOGIN"
	MESSAGE_TYPE_PING   = "PING"
	MESSAGE_TYPE_EFFECT = "EFFECT"
	// 返回类型
	MESSAGE_RESP_OK     = "OK"
	MESSAGE_RESP_ERROR  = "ERROR"
	MESSAGE_RESP_PONG   = "PONG"
	MESSAGE_RESP_EFFECT = "EFFECT"
	MESSAGE_ERROR_TYPE  = "未知type 'LOGIN'|'EFFECT'|'PING'"
	MESSAGE_ERROR_FORM  = "数据格式不对.示例: ws.send(JSON.stringify({type:'LOGIN',data:'user-token'}))"
	MESSAGE_ERROR_TOKEN = "token无效"
)

func NewClient(conn *websocket.Conn, room_id int, on_close func(*Client)) *Client {
	client := new(Client)
	client.Conn = conn
	client.OnClose = on_close
	client.RoomId = room_id
	client.Sender = make(chan interface{}, 64)
	return client
}

type MessageClient struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
type MessageServer struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	Conn    *websocket.Conn
	User    *model.User
	RoomId  int
	OnClose func(*Client)
	Sender  chan interface{}
}

func (this *Client) Close() {
	this.Conn.Close()
	this.OnClose(this)
}

func (this *Client) ReadPump() {
	defer func() {
		fmt.Println("read-pump-ended")
		close(this.Sender)
		this.Close()
	}()
	fmt.Println("read-pump...")
	this.Conn.SetReadLimit(maxMessageSize)
	this.Conn.SetReadDeadline(time.Now().Add(pongWait))
	this.Conn.SetPongHandler(func(string) error {
		fmt.Println("延长pong-wait")
		this.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		// 客户端断开 主动发送的消息这里都能接收到 断开msg_type==1
		// 有err就需要break 进入close程序 重复从关闭的conn读值将panic
		msg_type, msg_content, err := this.Conn.ReadMessage()
		fmt.Println("客户端消息", msg_type, string(msg_content))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("客户端主动断开")
			}
			fmt.Printf("已经断开: %v", err)
			break
		}
		this.HandleClientMessage(msg_content)
	}
}

func (this *Client) WritePump() {
	pingTick := time.NewTicker(pingPeriod)
	defer func() {
		fmt.Println("write-pump-ended")
		pingTick.Stop()
		this.Close()
	}()
	fmt.Println("write-pump...")
	for {
		select {
		case message, ok := <-this.Sender:
			if !ok {
				fmt.Println("通道已关闭")
				this.Conn.WriteMessage(websocket.CloseMessage, []byte(`{"data":"close"}`))
				return
			}
			fmt.Printf("comment: %+v\n", message.(model.Comment).Content)
			if err := this.SendMessage(MESSAGE_RESP_OK, message); err != nil {
				return
			}
		case <-pingTick.C:
			fmt.Println("ping-tick")
			this.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := this.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("ping-tick err:", err.Error())
				return
			}
		}
	}
}

func (this *Client) HandleClientMessage(msg_string []byte) {
	message, err := parseMessageClient(string(msg_string))
	if err != nil {
		this.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		this.Conn.WriteJSON(MessageServer{
			Type: MESSAGE_RESP_ERROR,
			Data: MESSAGE_ERROR_FORM,
		})
	}
	switch message.Type {
	case MESSAGE_TYPE_LOGIN:
		if user, err := parseLoginToken(message.Data); err != nil {
			this.SendMessage(MESSAGE_RESP_ERROR, MESSAGE_ERROR_TOKEN)
			break
		} else {
			this.User = user
			this.SendMessage(MESSAGE_RESP_OK, nil)
		}
	case MESSAGE_TYPE_PING:
		this.SendMessage(MESSAGE_RESP_PONG, nil)
	case MESSAGE_TYPE_EFFECT:
		fmt.Println("发送讨论组特效.暂不处理")
	default:
		this.SendMessage(MESSAGE_RESP_ERROR, MESSAGE_ERROR_TYPE)
	}
}

func (this *Client) SendMessage(msg_type string, msg_data interface{}) error {
	this.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := this.Conn.WriteJSON(MessageServer{
		Data: msg_data,
		Type: msg_type,
	})
	if err != nil {
		fmt.Println("推送消息失败 err:", err.Error())
	}
	return err
}

func parseMessageClient(msg_string string) (*MessageClient, error) {
	message := new(MessageClient)
	if err := json.Unmarshal([]byte(msg_string), message); err != nil {
		return nil, err
	}
	return message, nil
}

func parseLoginToken(token_string string) (*model.User, error) {
	tokenBody, err := auth.ParseToken(token_string)
	if err != nil {
		return nil, err
	}
	return &tokenBody.UserInfo, nil
}
