package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/google/uuid"
)

type BaseAgent[T IAgent[T]] struct {
	IExposedServerFunctions[T]
	id                      uuid.UUID
	messageLimiterSemaphore chan struct{}
}

func (a *BaseAgent[T]) GetID() uuid.UUID {
	return a.id
}

func CreateBaseAgent[T IAgent[T]](serv IExposedServerFunctions[T]) *BaseAgent[T] {
	return &BaseAgent[T]{
		IExposedServerFunctions: serv,
		id:                      uuid.New(),
		messageLimiterSemaphore: make(chan struct{}, serv.GetAgentMessagingBandwidth()),
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

func (a *BaseAgent[T]) SendMessage(msg message.IMessage[T], recipient uuid.UUID) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	select {
	case a.messageLimiterSemaphore <- struct{}{}:
		go func() {
			a.DeliverMessage(msg, recipient)
			<-a.messageLimiterSemaphore
		}()
	default:
	}
}

func (a *BaseAgent[T]) SendSynchronousMessage(msg message.IMessage[T], recipient uuid.UUID) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	a.DeliverMessage(msg, recipient)
}

func (agent *BaseAgent[T]) BroadcastMessage(msg message.IMessage[T]) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	for id := range agent.ViewAgentIdSet() {
		if id == msg.GetSender() {
			continue
		}
		agent.SendMessage(msg, id)
	}
}
