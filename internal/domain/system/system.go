package system

import "live-chat-kafka/internal/message_queue/types"

type ServerInfo struct {
	IP        string `json:"ip"`
	Available bool   `json:"available"`
}

func NewServerInfo(ip string, available bool) *ServerInfo {
	return &ServerInfo{
		IP:        ip,
		Available: available,
	}
}

func (s *ServerInfo) ConvertRedisData() map[string]interface{} {
	return map[string]interface{}{
		"ip":        s.IP,
		"available": s.Available,
	}
}

type UseCase interface {
	GetServerList() ([]ServerInfo, error)
	LoopSubKafka(timeoutMs int) (*types.Message, error)
	RegisterSubTopic(topic string) error
	SetChatServerInfo(ip string, available bool) error
	PublishServerStatusEvent(addr string, status bool)
}

type Repository interface {
	GetAvailableServerList() ([]ServerInfo, error)
	SetChatServerInfo(ip string, available bool) error
}

type PubSub interface {
	RegisterSubTopic(topic string) error
	Poll(timeoutMs int) types.Event
	PublishEvent(topic string, data []byte) (types.Event, error)
}
