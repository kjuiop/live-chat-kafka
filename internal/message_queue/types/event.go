package types

type Event interface {
	IsMessage() bool
	IsError() bool
}

type Message struct {
	Value []byte
}

func (m *Message) IsMessage() bool { return true }
func (m *Message) IsError() bool   { return false }

type Error struct {
	Error error
}

func (e *Error) IsMessage() bool { return false }
func (e *Error) IsError() bool   { return true }
