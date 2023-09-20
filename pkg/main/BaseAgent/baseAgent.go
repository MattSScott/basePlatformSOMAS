package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"
	"github.com/google/uuid"
)

type BaseAgent[T IAgent[T]] struct {
	id      uuid.UUID
	network []T
}

func (ba *BaseAgent[T]) GetID() uuid.UUID {
	return ba.id
}

func (ba *BaseAgent[T]) UpdateAgentInternalState() {
}

func NewAgent[T IAgent[T]]() *BaseAgent[T] {
	return &BaseAgent[T]{
		id: uuid.New(),
	}
}

func (ba *BaseAgent[T]) GetAllMessages() []message.IMessage[T] {
	return []message.IMessage[T]{}
}

func (ba *BaseAgent[T]) GetNetwork() []T {
	return ba.network
}

func (ba *BaseAgent[T]) SetNetwork(newNetwork []T) {
	ba.network = newNetwork
}
