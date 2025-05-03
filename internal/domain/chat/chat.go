package chat

import (
	"context"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	GetChatRoom(ctx context.Context, roomId string) (*Room, error)
	ServeWsByMemory(ctx context.Context, socket *websocket.Conn, chatRoom *Room, userId string) error
	ServeWs(ctx context.Context, socket *websocket.Conn, chatRoom *Room, userId string) error
}

type PubSub interface {
	SubscribeTopic(ctx context.Context, topic string, handler func(msg *Message)) error
	PublishMessage(roomId string, message *Message) error
}
