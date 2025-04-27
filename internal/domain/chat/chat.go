package chat

import (
	"context"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	GetChatRoom(ctx context.Context, roomId string) (*Room, error)
	ServeWs(ctx context.Context, socket *websocket.Conn, chatRoom *Room, userId string) error
}
