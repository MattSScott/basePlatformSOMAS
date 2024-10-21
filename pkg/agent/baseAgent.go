package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/internal/diagnosticsEngine"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/google/uuid"
)

type BaseAgent[T IAgent[T]] struct {
	IExposedServerFunctions[T]
	id                      uuid.UUID
	messageLimiterSemaphore chan struct{}
	diagnosticsEngine       diagnosticsEngine.IDiagnosticsEngine
}

func (a *BaseAgent[T]) GetID() uuid.UUID {
	return a.id
}

func CreateBaseAgent[T IAgent[T]](serv IExposedServerFunctions[T]) *BaseAgent[T] {
	if serv == nil {
		panic("Nil interface passed to CreateBaseAgent. Please pass an instance of IExposedServerFunctions")
	}
	return &BaseAgent[T]{
		IExposedServerFunctions: serv,
		id:                      uuid.New(),
		messageLimiterSemaphore: make(chan struct{}, serv.GetAgentMessagingBandwidth()),
		diagnosticsEngine:       serv.GetDiagnosticEngine(),
	}
}

func (a *BaseAgent[T]) CreateBaseMessage() message.BaseMessage {
	return message.BaseMessage{Sender: a.GetID()}
}

func (a *BaseAgent[T]) NotifyAgentFinishedMessaging() {
	go a.AgentStoppedTalking(a.id)
}

func (a *BaseAgent[T]) SendMessage(msg message.IMessage[T], recipient uuid.UUID) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	status := false
	select {
	case a.messageLimiterSemaphore <- struct{}{}:
		go func() {
			a.DeliverMessage(msg, recipient)
			<-a.messageLimiterSemaphore
		}()
		status = true
	default:
	}
	a.diagnosticsEngine.ReportSendMessageStatus(status)
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

func (agent *BaseAgent[T]) BroadcastSynchronousMessage(msg message.IMessage[T]) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	for id := range agent.ViewAgentIdSet() {
		if id == msg.GetSender() {
			continue
		}
		agent.SendSynchronousMessage(msg, id)
	}
}
