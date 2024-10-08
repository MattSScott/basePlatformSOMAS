package message

import (
	"github.com/google/uuid"
)

// base interface structure used for message - can be composed for more complex message structures
type IMessage[T any] interface {
	// returns the sender of a message
	GetSender() uuid.UUID
	// Sets the sender struct member of a message
	SetSender(uuid.UUID)
	// calls the appropriate message handler method on the receiving agent
	InvokeMessageHandler(T)
}

// new message types can extend this
type BaseMessage struct {
	sender uuid.UUID
}

func (bm *BaseMessage) GetSender() uuid.UUID {
	return bm.sender
}

func (bm *BaseMessage) SetSender(id uuid.UUID) {
	bm.sender = id
}
