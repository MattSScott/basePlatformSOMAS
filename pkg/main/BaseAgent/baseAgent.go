package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"
	"github.com/google/uuid"
)

type BaseAgent[T IAgent[T]] struct {
	id uuid.UUID
}

func (ba *BaseAgent[T]) GetID() uuid.UUID {
	return ba.id
}

func (ba *BaseAgent[T]) UpdateAgentInternalState() {}

func NewBaseAgent[T IAgent[T]]() *BaseAgent[T] {
	return &BaseAgent[T]{
		id: uuid.New(),
	}
}

func (ba *BaseAgent[T]) GetAllMessages(listOfAgents []T) []message.IMessage[T] {

	return []message.IMessage[T]{}
}
