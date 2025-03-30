package repository

import (
	"encoding/json"
	"fmt"
	"live-chat-kafka/internal/database"
	"live-chat-kafka/internal/domain/system"
	"live-chat-kafka/internal/models"
)

type systemRepository struct {
	db database.Client
}

func NewSystemRepository(db database.Client) system.Repository {
	return &systemRepository{
		db: db,
	}
}

func (s *systemRepository) GetAvailableServerList() ([]system.ServerInfo, error) {

	data, err := s.db.GetAvailableServerList()
	if err != nil {
		return nil, models.GetCustomErr(models.ErrNotFoundServerInfo)
	}

	if len(data) == 0 {
		return nil, nil
	}

	list := make([]system.ServerInfo, 0, len(data))
	for _, val := range data {
		var serverInfo system.ServerInfo
		if err := json.Unmarshal([]byte(val), &serverInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal server info: %w", err)
		}
		list = append(list, serverInfo)
	}

	return list, nil
}

func (s *systemRepository) SetChatServerInfo(ip string, available bool) error {

	if err := s.db.SaveChatServerInfo(ip, available); err != nil {
		return err
	}

	return nil
}
