package infra

import (
	"fmt"

	"github.com/google/uuid"
)

// // base interface structure used for message passing - can be composed for more complex message structures
type IAgentMessenger[T any] interface {
	// produces a list of messages (of any common interface) that an agent wishes to pass
	GetAllMessages(listOfAgents []T) []IMessage[T]
}

// base interface structure used for message - can be composed for more complex message structures
type IMessage[T any] interface {
	//MessagingProtocol[T]
	// returns the sender of a message
	GetSender() uuid.UUID
	// returns the list of agents that the message should be passed to
	//GetRecipients() []T
	// calls the appropriate messsage handler method on the receiving agent
	InvokeMessageHandler(uuid.UUID)
	Print()
}



// new message types can extend this
type BaseMessage[T IAgentMessenger[T]] struct {
	sender     T

}

// create read-only message instance
func CreateMessage[T IAgentMessenger[T]](sender T, recipients []T) BaseMessage[T] {
	return BaseMessage[T]{
		sender:     sender,
	}
}

func (bm BaseMessage[T]) Print() {
	fmt.Printf("message received from %s\n", bm.sender)
}

func (bm BaseMessage[T]) GetSender() T {
	return bm.sender
}


func (bm BaseMessage[T]) InvokeMessageHandler(agent T) {
}




