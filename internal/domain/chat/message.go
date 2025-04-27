package chat

type Message struct {
	SendUserId string
	Message    string
	Time       int64
	Method     string
}
