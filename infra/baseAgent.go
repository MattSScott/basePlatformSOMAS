package infra

import (
	"github.com/google/uuid"
)

type BaseAgent[T IAgent[T]] struct {
	IExposedServerFunctions[T]
	id uuid.UUID
}

func (ba *BaseAgent[T]) GetID() uuid.UUID {
	return ba.id
}

func CreateBaseAgent[T IAgent[T]](serv IExposedServerFunctions[T]) *BaseAgent[T] {
	return &BaseAgent[T]{
		IExposedServerFunctions: serv,
		id:                      uuid.New(),
	}
}

func (a *BaseAgent[T]) UpdateAgentInternalState() {}

func (a *BaseAgent[T]) NotifyAgentFinishedMessaging() {
	a.agentStoppedTalking(a.id)
}

func (a *BaseAgent[T]) RunSynchronousMessaging() {}
