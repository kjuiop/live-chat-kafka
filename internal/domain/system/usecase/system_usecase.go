package usecase

import (
	"errors"
	"live-chat-kafka/internal/domain/system"
	"live-chat-kafka/internal/models"
	"log"
)

type systemUseCase struct {
	systemRepo    system.Repository
	systemPubSub  system.PubSub
	avgServerList map[string]bool
}

func NewSystemUseCase(systemRepo system.Repository, systemPubSub system.PubSub) system.UseCase {

	s := &systemUseCase{
		systemRepo:    systemRepo,
		systemPubSub:  systemPubSub,
		avgServerList: make(map[string]bool),
	}

	if err := s.setServerInfo(); err != nil {
		log.Fatalf("failed register server info, err : %v", err)
	}

	return s
}

func (s *systemUseCase) GetServerList() ([]system.ServerInfo, error) {

	if len(s.avgServerList) == 0 {
		return nil, nil
	}

	var res []system.ServerInfo

	for ip, available := range s.avgServerList {
		if len(ip) > 0 && available {
			server := system.ServerInfo{
				IP:        ip,
				Available: available,
			}
			res = append(res, server)
		}
	}

	return res, nil
}

func (s *systemUseCase) GetAvailableServerList() ([]system.ServerInfo, error) {
	return s.systemRepo.GetAvailableServerList()
}

func (s *systemUseCase) setServerInfo() error {

	serverList, err := s.GetAvailableServerList()
	if errors.Is(err, models.GetCustomErr(models.ErrNotFoundServerInfo)) {
		return nil
	} else if err != nil {
		return err
	}

	if len(serverList) == 0 {
		return nil
	}

	for _, server := range serverList {
		s.avgServerList[server.IP] = true
	}

	return nil
}
