package usecase

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"live-chat-kafka/internal/domain/chat"
	"live-chat-kafka/internal/domain/chat/types"
	"live-chat-kafka/internal/domain/room"
	"net/http"
	"sync"
	"time"
)

var crMutex = &sync.RWMutex{}

type chatUseCase struct {
	roomUseCase    room.UseCase
	contextTimeout time.Duration
	hub            map[string]*chat.Room
	upgrader       *websocket.Upgrader
	chatPubSub     chat.PubSub
}

func NewChatUseCase(roomUseCase room.UseCase, timeout time.Duration, chatPubSub chat.PubSub) chat.UseCase {
	return &chatUseCase{
		roomUseCase:    roomUseCase,
		contextTimeout: timeout,
		hub:            make(map[string]*chat.Room),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  types.SocketBufferSize,
			WriteBufferSize: types.MessageBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		chatPubSub: chatPubSub,
	}
}

func (cu *chatUseCase) ServeWs(ctx context.Context, socket *websocket.Conn, chatRoom *chat.Room, userId string) error {

	client := chat.NewClient(socket, chatRoom, userId)

	chatRoom.Join <- client
	defer func() {
		chatRoom.Leave <- client
	}()

	go client.Write(ctx)

	client.Read(ctx)

	return nil
}

func (cu *chatUseCase) ServeWsByMemory(ctx context.Context, socket *websocket.Conn, chatRoom *chat.Room, userId string) error {

	client := chat.NewClient(socket, chatRoom, userId)

	chatRoom.Join <- client
	defer func() {
		chatRoom.Leave <- client
	}()

	go client.Write(ctx)

	client.Read(ctx)

	return nil
}

func (cu *chatUseCase) GetChatRoom(ctx context.Context, roomId string) (*chat.Room, error) {
	crMutex.Lock()
	defer func() {
		crMutex.Unlock()
	}()

	if _, ok := cu.hub[roomId]; !ok {
		roomInfo, err := cu.roomUseCase.GetChatRoomById(ctx, roomId)
		if err != nil {
			return nil, fmt.Errorf("not found chat room, key : %s, err : %w", roomId, err)
		}

		chatRoom, err := chat.NewChatRoom(ctx, roomInfo, cu.chatPubSub)
		if err != nil {
			return nil, fmt.Errorf("not found chat room, key : %s, err : %w", roomId, err)
		}

		cu.hub[roomId] = chatRoom
	}

	return cu.hub[roomId], nil
}
