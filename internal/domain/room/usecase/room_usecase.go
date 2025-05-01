package usecase

import (
	"context"
	"live-chat-kafka/api/form"
	"live-chat-kafka/internal/domain/room"
	"time"
)

type roomUseCase struct {
	roomRepo       room.Repository
	contextTimeout time.Duration
	roomPubSub     room.PubSub
}

func NewRoomUseCase(roomRepo room.Repository, timeout time.Duration, roomPubSub room.PubSub) room.UseCase {
	return &roomUseCase{
		roomRepo:       roomRepo,
		contextTimeout: timeout,
		roomPubSub:     roomPubSub,
	}
}

func (r *roomUseCase) CreateChatRoom(c context.Context, room room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.Create(ctx, room); err != nil {
		return err
	}

	if err := r.roomPubSub.CreateChatRoom(c, room.RoomId); err != nil {
		return err
	}

	return nil
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

	if err := r.roomPubSub.DeleteChatRoom(ctx, roomId); err != nil {
		return err
	}

	return nil
}

func (r *roomUseCase) GetChatRoomId(c context.Context, req form.RoomIdRequest) (*room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	roomInfo, err := r.roomRepo.GetRoomMap(ctx, req.ChannelKey, req.BroadCastKey)
	if err != nil {
		return nil, err
	}

	return roomInfo, nil
}
