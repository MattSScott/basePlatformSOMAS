package message

import (
	"fmt"

	"github.com/google/uuid"
)

// base interface structure used for message - can be composed for more complex message structures
type IMessage[T any] interface {
	// returns the sender of a message
	GetSender() uuid.UUID
	SetSender(uuid.UUID)
	// calls the appropriate messsage handler method on the receiving agent
	InvokeMessageHandler(T)
	// prints message to console
	Print()
}

// new message types can extend this
type BaseMessage struct {
	sender uuid.UUID
}

func (bm *BaseMessage) Print() {
	fmt.Printf("message received from %s\n", bm.sender)
}

func (bm *BaseMessage) GetSender() uuid.UUID {
	return bm.sender
}

func (bm *BaseMessage) SetSender(id uuid.UUID) {
	bm.sender = id
}
