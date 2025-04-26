package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"live-chat-kafka/internal/database"
	"live-chat-kafka/internal/domain/room"
	"time"
)

const (
	LiveChatServerRoom    = "live-chat-server-room"
	LiveChatServerRoomMap = "live-chat-server-room-map"
)

const (
	RoomExpire = time.Duration(7) * 24 * time.Hour
)

type roomRepository struct {
	db database.Client
}

func NewRoomRepository(db database.Client) room.Repository {
	return &roomRepository{
		db: db,
	}
}

func (r roomRepository) Create(ctx context.Context, data room.RoomInfo) error {

	if err := r.db.Set(ctx, convertRoomKey(data.RoomId), data.ConvertRedisData(), RoomExpire); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	return nil
}

func (r roomRepository) SetRoomMap(ctx context.Context, data room.RoomInfo) error {

	jData, err := json.Marshal(data.ConvertRedisData())
	if err != nil {
		return fmt.Errorf("set room map json encoding fail, err : %w", err)
	}

	if err := r.db.Set(ctx, generateRoomMapKey(data.ChannelKey, data.BroadcastKey), string(jData), RoomExpire); err != nil {
		return fmt.Errorf("set room map err : %w", err)
	}

	return nil
}

func convertRoomKey(roomId string) string {
	return fmt.Sprintf("%s_%s", LiveChatServerRoom, roomId)
}

func generateRoomMapKey(channelKey, broadcastKey string) string {
	return fmt.Sprintf("%s_%s_%s", LiveChatServerRoomMap, channelKey, broadcastKey)
}
