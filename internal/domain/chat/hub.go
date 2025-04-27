package chat

import (
	"live-chat-kafka/internal/domain/room"
	"log/slog"
)

type Room struct {
	RoomId string `json:"RoomId"`
	Alive  bool   `json:"Alive"`

	Forward chan *Message // 수신되는 메시지를 보관하는 값
	// 들어오는 메시지를 다른 클라이언트에게 전송을 합니다.

	Join  chan *Client // Socket 이 연결되는 경우에 적용
	Leave chan *Client // Socket 이 끊어지는 경우에 대해서 적용

	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장
}

func NewChatRoom(roomInfo *room.RoomInfo) *Room {
	chatRoom := &Room{
		RoomId:  roomInfo.RoomId,
		Alive:   true,
		Forward: make(chan *Message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
	go chatRoom.chatInit()
	return chatRoom
}

func (r *Room) chatInit() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
			slog.Info("join", "client_id", client.UserID, "room_id", client.Room.RoomId)
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
			slog.Info("leave", "client_id", client.UserID, "room_id", client.Room.RoomId)
		case msg := <-r.Forward:
			// msg type 에 따른 분기 처리 가능
			switch msg.Method {
			case "chat":
				r.broadcastChat(msg)
				slog.Info("chat", "send_user_id", msg.SendUserId, "chat_message", msg.Message)
			default:
			}
		}
	}
}

func (r *Room) broadcastChat(msg *Message) {
	for client := range r.Clients {
		client.Send <- msg
	}
}
