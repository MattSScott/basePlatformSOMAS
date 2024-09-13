package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/google/uuid"
)

type BaseAgent[T IAgent[T]] struct {
	IExposedServerFunctions[T]
	id uuid.UUID
}

func (a *BaseAgent[T]) GetID() uuid.UUID {
	return a.id
}

func CreateBaseAgent[T IAgent[T]](serv IExposedServerFunctions[T]) *BaseAgent[T] {
	return &BaseAgent[T]{
		IExposedServerFunctions: serv,
		id:                      uuid.New(),
	}
}

func (a *BaseAgent[T]) CreateBaseMessage() message.BaseMessage {
	msg := message.BaseMessage{}
	msg.SetSender(a.GetID())
	return msg
}

func (a *BaseAgent[T]) UpdateAgentInternalState() {}

func (a *BaseAgent[T]) NotifyAgentFinishedMessaging() {
	go a.AgentStoppedTalking(a.id)
}

func (a *BaseAgent[T]) RunSynchronousMessaging() {}
