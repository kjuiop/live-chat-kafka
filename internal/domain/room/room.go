package room

import (
	"context"
	"fmt"
	"live-chat-kafka/api/form"
	"live-chat-kafka/utils"
	"math/rand"
	"strings"
	"time"
)

type RoomInfo struct {
	RoomId       string `json:"room_id"`
	CustomerId   string `json:"customer_id"`
	ChannelKey   string `json:"channel_key"`
	BroadcastKey string `json:"broadcast_key"`
	CreatedAt    int64  `json:"created_at"`
}

func NewRoomInfo(req form.RoomRequest, prefix string) *RoomInfo {
	return &RoomInfo{
		RoomId:       fmt.Sprintf("%s-%s", getChatPrefix(prefix), utils.GenUUID()),
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadcastKey: req.BroadCastKey,
		CreatedAt:    time.Now().Unix(),
	}
}

func UpdateRoomInfo(req form.RoomRequest, roomId string) *RoomInfo {
	return &RoomInfo{
		RoomId:       roomId,
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadcastKey: req.BroadCastKey,
		CreatedAt:    time.Now().Unix(),
	}
}

func (r *RoomInfo) ConvertRedisData() map[string]interface{} {
	return map[string]interface{}{
		"room_id":       r.RoomId,
		"customer_id":   r.CustomerId,
		"channel_key":   r.ChannelKey,
		"broadcast_key": r.BroadcastKey,
		"created_at":    r.CreatedAt,
	}
}

func getChatPrefix(prefix string) string {
	array := strings.Split(prefix, ",")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[rand.Intn(len(array))]
}

type UseCase interface {
	CreateChatRoom(ctx context.Context, room RoomInfo) error
	GetChatRoomById(ctx context.Context, roomId string) (*RoomInfo, error)
	RegisterRoomId(ctx context.Context, room RoomInfo) error
	CheckExistRoomId(ctx context.Context, roomId string) (bool, error)
	UpdateChatRoom(ctx context.Context, roomId string, room RoomInfo) (*RoomInfo, error)
}

type Repository interface {
	Create(ctx context.Context, data RoomInfo) error
	Fetch(ctx context.Context, roomId string) (*RoomInfo, error)
	Exists(ctx context.Context, roomId string) (bool, error)
	Update(ctx context.Context, roomId string, data RoomInfo) error
	SetRoomMap(ctx context.Context, data RoomInfo) error
}
