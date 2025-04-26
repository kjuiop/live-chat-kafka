package usecase

import "live-chat-kafka/internal/domain/room"

type roomUseCase struct {
	roomRepo room.Repository
}

func NewRoomUseCase(roomRepo room.Repository) room.UseCase {
	return &roomUseCase{
		roomRepo: roomRepo,
	}
}
