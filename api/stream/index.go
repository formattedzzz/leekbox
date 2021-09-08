package stream

import (
	"fmt"
	"leekbox/model"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	DEV_MODE = os.Getenv("LEEKBOX_ENV")
)

type StreamAPI struct {
	RoomClient  map[int][]*Client
	lock        sync.RWMutex
	PingPeriod  time.Duration
	PongTimeout time.Duration
	Upgrader    *websocket.Upgrader
}

func New(ping time.Duration, pong time.Duration, allowed_origins []string) *StreamAPI {
	return &StreamAPI{
		RoomClient:  make(map[int][]*Client),
		PingPeriod:  ping,
		PongTimeout: ping + pong,
		Upgrader:    newUpgrader(allowed_origins),
	}
}

func (this *StreamAPI) RemoveClient(client *Client) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if client_list, ok := this.RoomClient[client.RoomId]; ok {
		for i, c := range client_list {
			if c == client {
				this.RoomClient[client.RoomId] = append(client_list[:i], client_list[i+1:]...)
			}
		}
		fmt.Println(fmt.Sprintf("removed %d", len(client_list)))
	}
}

func (this *StreamAPI) RegisterClient(client *Client) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.RoomClient[client.RoomId] = append(this.RoomClient[client.RoomId], client)
}

func (this *StreamAPI) PushRoomComment(comment *model.Comment, userInfo model.User) {
	if client_list, ok := this.RoomClient[comment.RoomId]; ok {
		for _, client := range client_list {
			client.Sender <- model.CommentWithUser{
				Comment: comment,
				User:    userInfo,
			}
		}
	}
}

// @Summary 讨论组发言流websocket
// @Param room_id query string true "讨论组ID"
// @Description 客户端接口示例
// @Description ws = new WebSocket("ws://host:port/api/stream?room_id=1")
// @Description ws.addEventListener("close", console.log)
// @Description ws.addEventListener("error", console.log)
// @Description ws.addEventListener("message", ev => console.log("message:", ev.data))
// @Description ws.addEventListener("open", ev => {
// @Description
// @Description   console.log("连接成功")
// @Description		连接鉴权
// @Description   ws.send(JSON.stringify({ data: "user-token", type: "LOGIN" }))
// @Description		发送临时消息
// @Description   ws.send(JSON.stringify({ data: "hello world!", type: "MESSAGE" }))
// @Description		发送心跳包
// @Description   ws.send(JSON.stringify({ data: "", type: "PING" }))
// @Description		发送临时特效代码
// @Description   ws.send(JSON.stringify({ data: "1", type: "EFFECT" }))
// @Description
// @Description })
// @Router /api/stream [get]
func (this *StreamAPI) Handler(ctx *gin.Context) {
	room_id := 0
	if id, err := strconv.Atoi(ctx.Query("room_id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Return(40000, nil, "讨论组ID无效"))
		return
	} else {
		room_id = id
	}
	conn, err := this.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Return(40000, nil, err.Error()))
		return
	}
	client := NewClient(this, conn, room_id, this.RemoveClient)
	this.RegisterClient(client)
	fmt.Printf("client: %v\n", *client)
	go client.ReadPump()
	go client.WritePump()
}

func newUpgrader(allowedWebSocketOrigins []string) *websocket.Upgrader {
	compiledAllowedOrigins := compileAllowedWebSocketOrigins(allowedWebSocketOrigins)
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if DEV_MODE != "PROD" {
				return true
			}
			return isAllowedOrigin(r, compiledAllowedOrigins)
		},
	}
}

func isAllowedOrigin(r *http.Request, allowedOrigins []*regexp.Regexp) bool {
	origin := r.Header.Get("origin")
	if origin == "" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if strings.EqualFold(u.Host, r.Host) {
		return true
	}
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin.Match([]byte(strings.ToLower(u.Hostname()))) {
			return true
		}
	}
	return false
}
func compileAllowedWebSocketOrigins(allowedOrigins []string) []*regexp.Regexp {
	var compiledAllowedOrigins []*regexp.Regexp
	for _, origin := range allowedOrigins {
		compiledAllowedOrigins = append(compiledAllowedOrigins, regexp.MustCompile(origin))
	}
	return compiledAllowedOrigins
}
