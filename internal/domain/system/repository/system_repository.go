package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"live-chat-kafka/internal/database"
	"live-chat-kafka/internal/domain/system"
	"live-chat-kafka/internal/models"
)

const (
	LiveChatServerInfo = "live-chat-server-info"
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

	data, err := s.db.HGetAll(context.TODO(), LiveChatServerInfo)
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

	data := system.NewServerInfo(ip, available)
	if err := s.db.HSet(context.TODO(), LiveChatServerInfo, data.IP, data.ConvertRedisData(), 0); err != nil {
		return err
	}

	return nil
}
