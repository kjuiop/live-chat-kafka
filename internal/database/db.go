package database

type Client interface {
	GetAvailableServerList() (map[string]string, error)
	SaveChatServerInfo(ip string, available bool) error
}
