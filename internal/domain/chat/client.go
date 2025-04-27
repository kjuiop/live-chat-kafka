package chat

import (
	"github.com/gorilla/websocket"
	"live-chat-kafka/internal/domain/chat/types"
	"log/slog"
	"time"
)

type Client struct {
	Send   chan *Message
	Room   *Room
	UserID string
	Socket *websocket.Conn
}

func NewClient(socket *websocket.Conn, r *Room, clientId string) *Client {
	return &Client{
		Socket: socket,
		Send:   make(chan *Message, types.MessageBufferSize),
		Room:   r,
		UserID: clientId,
	}
}

func (c *Client) Write() {
	defer func() {
		if err := c.Socket.Close(); err != nil {
			slog.Error("failed socket connection close, client_id : %s, err : %s", c.UserID, err.Error())
		}
	}()
	// client 가 메시지를 전송하는 함수

	for msg := range c.Send {
		if err := c.Socket.WriteJSON(msg); err != nil {
			slog.Error("failed write message, client_id : %s, err : %s", c.UserID, err.Error())
		}
	}
}

func (c *Client) Read() {
	defer func() {
		if err := c.Socket.Close(); err != nil {
			slog.Error("failed socket connection close, client_id : %s, err : %s", c.UserID, err.Error())
		}
	}()

	// client 가 메시지를 읽는 함수
	for {
		var msg *Message
		if err := c.Socket.ReadJSON(&msg); err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				break
			}

			slog.Error("failed read message, client_id : %s, err : %s", c.UserID, err.Error())
			continue
		}
		msg.Time = time.Now().Unix()
		msg.SendUserId = c.UserID

		c.Room.Forward <- msg
	}
}
