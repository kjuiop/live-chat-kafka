package chat

import (
	"context"
	"github.com/gorilla/websocket"
	"live-chat-kafka/internal/domain/chat/types"
	"live-chat-kafka/internal/message_queue"
	"log/slog"
	"time"
)

type Client struct {
	Send     chan *Message
	Room     *Room
	UserID   string
	Socket   *websocket.Conn
	mq       message_queue.Client
	isClosed bool
}

func NewClient(socket *websocket.Conn, r *Room, clientId string) *Client {

	client := &Client{
		Socket: socket,
		Send:   make(chan *Message, types.MessageBufferSize),
		Room:   r,
		UserID: clientId,
	}

	return client
}

func (c *Client) Write(ctx context.Context) {
	defer func() {
		if !c.isClosed && c.Socket != nil {
			if err := c.Socket.Close(); err != nil {
				slog.Error("failed socket connection close", "client_id", c.UserID, "err", err)
			}
			c.isClosed = true
		}
	}()
	// client 가 메시지를 전송하는 함수
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.Send:
			if !ok {
				return
			}
			if err := c.Socket.WriteJSON(msg); err != nil {
				slog.Error("failed write message", "client_id", c.UserID, "err", err)
				return
			}
		}
	}
}

func (c *Client) Read(ctx context.Context) {
	defer func() {
		if !c.isClosed && c.Socket != nil {
			if err := c.Socket.Close(); err != nil {
				slog.Error("failed socket connection close", "client_id", c.UserID, "err", err)
			}
			c.isClosed = true
		}
	}()

	// client 가 메시지를 읽는 함수
ReadLoop:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg *Message
			if err := c.Socket.ReadJSON(&msg); err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
					break ReadLoop
				}

				slog.Error("failed read message", "client_id", c.UserID, "err", err)
				continue
			}
			msg.Time = time.Now().Unix()
			msg.SendUserId = c.UserID

			if err := c.Room.chatPubSub.PublishMessage(c.Room.RoomId, msg); err != nil {
				slog.Error("failed to publish chat message", "room_id", c.Room.RoomId, "user_id", c.UserID, "err", err)
			}
		}
	}
}
