package repository

import (
	"live-chat-kafka/internal/database"
	"live-chat-kafka/internal/domain/room"
)

type roomRepository struct {
	db database.Client
}

func NewRoomRepository(db database.Client) room.Repository {
	return &roomRepository{
		db: db,
	}
}
