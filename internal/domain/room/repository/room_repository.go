package repository

import (
	"context"
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

	if err := r.db.Set(ctx, generateRoomMapKey(data.ChannelKey, data.BroadcastKey), data.ConvertRedisData(), RoomExpire); err != nil {
		return fmt.Errorf("set room map err : %w", err)
	}

	return nil
}

func (r roomRepository) Fetch(ctx context.Context, roomId string) (*room.RoomInfo, error) {

	roomInfo := &room.RoomInfo{}
	if err := r.db.Get(ctx, convertRoomKey(roomId), roomInfo); err != nil {
		return nil, fmt.Errorf("get chat room info err : %w", err)
	}

	return roomInfo, nil
}

func (r roomRepository) Exists(ctx context.Context, roomId string) (bool, error) {

	isExist, err := r.db.Exists(ctx, convertRoomKey(roomId))
	if err != nil {
		return false, fmt.Errorf("fail redis cmd exist err : %w", err)
	}

	return isExist, nil
}

func (r roomRepository) Update(ctx context.Context, roomId string, data room.RoomInfo) error {

	if err := r.db.Set(ctx, convertRoomKey(roomId), data.ConvertRedisData(), RoomExpire); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	return nil
}

func (r roomRepository) Delete(ctx context.Context, roomId string) error {

	if err := r.db.DelByKey(ctx, convertRoomKey(roomId)); err != nil {
		return err
	}

	return nil
}

func (r roomRepository) GetRoomMap(ctx context.Context, channelKey, broadcastKey string) (*room.RoomInfo, error) {

	roomInfo := &room.RoomInfo{}
	if err := r.db.Get(ctx, generateRoomMapKey(channelKey, broadcastKey), roomInfo); err != nil {
		return nil, fmt.Errorf("get chat room map err : %w", err)
	}

	return roomInfo, nil
}

func convertRoomKey(roomId string) string {
	return fmt.Sprintf("%s_%s", LiveChatServerRoom, roomId)
}

func generateRoomMapKey(channelKey, broadcastKey string) string {
	return fmt.Sprintf("%s_%s_%s", LiveChatServerRoomMap, channelKey, broadcastKey)
}
