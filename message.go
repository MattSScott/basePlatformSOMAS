package basePlatformSOMAS

import (
	"fmt"

	"github.com/google/uuid"
)

// // base interface structure used for message passing - can be composed for more complex message structures
// type IAgentMessenger interface {
// 	// produces a list of messages (of any common interface) that an agent wishes to pass
// 	GetAllMessages(listOfAgents []T) []IMessage[T]
// }

// base interface structure used for message - can be composed for more complex message structures
type IMessage[T any] interface {
	// returns the sender of a message
	GetSender() uuid.UUID
	// calls the appropriate messsage handler method on the receiving agent
	InvokeMessageHandler(T)
	// prints message to console
	Print()
}

// new message types can extend this
type BaseMessage struct {
	sender uuid.UUID
}

// create read-only message instance


func (bm BaseMessage) Print() {
	fmt.Printf("message received from %s\n", bm.sender)
}

func (bm BaseMessage) GetSender() uuid.UUID {
	return bm.sender
}

