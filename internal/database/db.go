package database

type Client interface {
	GetAvailableServerList() (map[string]string, error)
	SaveChatServerInfo(key string, data map[string]interface{}) error
}
