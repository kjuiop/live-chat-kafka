package database

type Client interface {
	GetAvailableServerList() (map[string]string, error)
}
