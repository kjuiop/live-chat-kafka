package usecase

import (
	"context"
	"live-chat-kafka/internal/domain/room"
	"time"
)

type roomUseCase struct {
	roomRepo       room.Repository
	contextTimeout time.Duration
}

func NewRoomUseCase(roomRepo room.Repository, timeout time.Duration) room.UseCase {
	return &roomUseCase{
		roomRepo:       roomRepo,
		contextTimeout: timeout,
	}
}

func (r roomUseCase) CreateChatRoom(c context.Context, room room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	return r.roomRepo.Create(ctx, room)
}

func (r roomUseCase) RegisterRoomId(c context.Context, roomInfo room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.SetRoomMap(ctx, roomInfo); err != nil {
		return err
	}

	return nil
}
