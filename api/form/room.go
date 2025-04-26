package form

type RoomRequest struct {
	CustomerId   string `json:"customer_id" binding:"required"`
	ChannelKey   string `json:"channel_key" binding:"required"`
	BroadCastKey string `json:"broadcast_key" binding:"required"`
}

type RoomResponse struct {
	RoomId       string `json:"room_id"`
	CustomerId   string `json:"customer_id,omitempty"`
	ChannelKey   string `json:"channel_key,omitempty"`
	BroadcastKey string `json:"broadcast_key,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
}
