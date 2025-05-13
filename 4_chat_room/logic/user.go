package logic

import (
	"context"
	"errors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strings"
	"time"
)

var globalUID uint32 = 0

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`
	Token          string        `json:"token"`

	conn *websocket.Conn

	isNew bool
}

// 系统用户，代表是系统主动发送的消息
var System = &User{}

func NewUser(conn *websocket.Conn, nickname, addr string) *User {
	return &User{
		conn:     conn,
		NickName: nickname,
		Addr:     addr,
	}
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)

	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			// 判定连接是否关闭了，正常关闭，不认为是错误
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			return err
		}
		// 内容发送到聊天室
		sendMsg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])

		// 解析 content，看是否是一条私信消息
		sendMsg.Content = strings.TrimSpace(sendMsg.Content)
		if strings.HasPrefix(sendMsg.Content, "@") {
			sendMsg.To = strings.SplitN(sendMsg.Content, " ", 2)[0][1:]
		}

		Broadcaster.Broadcast(sendMsg)
	}
}

func (u *User) SendMessageByToken(ctx context.Context, token string) error {
	return nil
}

func (u *User) CloseMessageChannel() {

}
