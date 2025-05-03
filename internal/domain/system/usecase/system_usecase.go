package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"live-chat-kafka/internal/domain/system"
	"live-chat-kafka/internal/message_queue/types"
	"live-chat-kafka/internal/models"
	"log"
	"log/slog"
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

	if err := s.initProcess(); err != nil {
		log.Fatalf("system usecase init process error: %v", err)
	}

	return s
}

func (s *systemUseCase) initProcess() error {

	if err := s.getServerInfoForMemory(); err != nil {
		return fmt.Errorf("get server info for memory error: %w", err)
	}

	if err := s.createChatServerTopic(); err != nil {
		return fmt.Errorf("create chat server topic error: %w", err)
	}

	return nil
}

func (s *systemUseCase) createChatServerTopic() error {

	if err := s.systemPubSub.CreateChatServerTopic(); err != nil {
		return err
	}

	return nil
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

func (s *systemUseCase) getServerInfoForMemory() error {

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

func (s *systemUseCase) LoopSubKafka(timeoutMs int) (*types.Message, error) {

	ev := s.systemPubSub.Poll(timeoutMs) // Polling 1000ms 동안 이벤트 대기

	if ev == nil {
		return nil, nil
	}

	if ev.IsError() {
		errorEvent := ev.(*types.Error)
		slog.Error("Failed to Polling event", "error", errorEvent.Error)
		return nil, fmt.Errorf("consumer event error, err : %v", errorEvent.Error)
	}

	if !ev.IsMessage() {
		return nil, fmt.Errorf("is not expected message, event : %v", ev)
	}

	message := ev.(*types.Message)

	var decoder system.ServerInfo
	if err := json.Unmarshal(message.Value, &decoder); err != nil {
		slog.Error("failed to decode event", "event_value", string(message.Value))
		return nil, err
	}

	if err := s.SetChatServerInfo(decoder.IP, decoder.Available); err != nil {
		slog.Error("failed update server info", "server_ip", decoder.IP, "available", decoder.Available)
		return nil, err
	}

	s.avgServerList[decoder.IP] = decoder.Available

	slog.Debug("update chat server info", "server_ip", decoder.IP, "available", decoder.Available, "avg_server_list", s.avgServerList)

	return &types.Message{Value: message.Value}, nil
}

func (s *systemUseCase) SetChatServerInfo(ip string, available bool) error {
	if err := s.systemRepo.SetChatServerInfo(ip, available); err != nil {
		return err
	}
	return nil
}

func (s *systemUseCase) SubscribeTopic(topic string) error {
	return s.systemPubSub.SubscribeTopic(topic)
}

func (s *systemUseCase) PublishServerStatusEvent(addr string, status bool) {

	serverInfo := system.ServerInfo{IP: addr, Available: status}

	bytes, err := json.Marshal(serverInfo)
	if err != nil {
		log.Fatalf("failed register server info, address : %s, err : %v", addr, err)
	}

	event, err := s.systemPubSub.PublishEvent("chat", bytes)
	if err != nil {
		log.Fatalf("failed publish server info, addr : %s, err : %v", addr, err.Error())
	}

	slog.Debug("success server info publish", "event", event)
}
