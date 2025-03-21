package system

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
}

type Repository interface {
	GetAvailableServerList() ([]ServerInfo, error)
}

type PubSub interface {
}
