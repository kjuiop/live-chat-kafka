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

func (r *roomUseCase) CreateChatRoom(c context.Context, room room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	return r.roomRepo.Create(ctx, room)
}

func (r *roomUseCase) RegisterRoomId(c context.Context, roomInfo room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.SetRoomMap(ctx, roomInfo); err != nil {
		return err
	}

	return nil
}

func (r *roomUseCase) GetChatRoomById(c context.Context, roomId string) (*room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	roomInfo, err := r.roomRepo.Fetch(ctx, roomId)
	if err != nil {
		return nil, err
	}

	return roomInfo, nil
}

func (r *roomUseCase) CheckExistRoomId(c context.Context, roomId string) (bool, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	isExist, err := r.roomRepo.Exists(ctx, roomId)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (r *roomUseCase) UpdateChatRoom(c context.Context, roomId string, roomInfo room.RoomInfo) (*room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.Update(ctx, roomId, roomInfo); err != nil {
		return nil, err
	}

	savedInfo, err := r.roomRepo.Fetch(c, roomId)
	if err != nil {
		return nil, err
	}

	return savedInfo, nil
}

func (r *roomUseCase) DeleteChatRoom(c context.Context, roomId string) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.Delete(ctx, roomId); err != nil {
		return err
	}
	return nil
}
