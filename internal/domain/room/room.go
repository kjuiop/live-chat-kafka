package room

import (
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

func getChatPrefix(prefix string) string {
	array := strings.Split(prefix, ",")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[rand.Intn(len(array))]
}

type UseCase interface {
}

type Repository interface {
}
