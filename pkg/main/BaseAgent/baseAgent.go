package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"
	"github.com/google/uuid"
)

type BaseAgent[T IAgent[T]] struct {
	id      uuid.UUID
	network map[uuid.UUID]T
}

func (ba *BaseAgent[T]) GetID() uuid.UUID {
	return ba.id
}

func (ba *BaseAgent[T]) UpdateAgentInternalState() {}

func NewAgent[T IAgent[T]]() *BaseAgent[T] {
	return &BaseAgent[T]{
		id:      uuid.New(),
		network: make(map[uuid.UUID]T),
	}
}

func DefaultAgent[T IAgent[T]]() IAgent[T] {
	return &BaseAgent[T]{
		id:      uuid.New(),
		network: make(map[uuid.UUID]T),
	}
}

func (ba *BaseAgent[T]) GetAllMessages(listOfAgents []T) []message.IMessage[T] {

	return []message.IMessage[T]{}
}

func (ba *BaseAgent[T]) GetNetwork() map[uuid.UUID]T {
	return ba.network
}

func (ba *BaseAgent[T]) AddAgentToNetwork(agent T) {
	ba.network[agent.GetID()] = agent
}

func (ba *BaseAgent[T]) RemoveAgentFromNetwork(agent T) {
	delete(ba.network, agent.GetID())
}
