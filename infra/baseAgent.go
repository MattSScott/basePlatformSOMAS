package infra

import (
	"fmt"
	"sync"

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

func (a *BaseAgent[T]) NotifyAgentInactive() {
	a.agentStoppedTalking(a.id)
}

func (a *BaseAgent[T]) RunSynchronousMessaging() {}

func (a *BaseAgent[T]) listenOnChannel(agentAgentchannel chan IMessage, serverAgentchannel chan ServerNotification, wait *sync.WaitGroup) {
	defer wait.Done()

	// checkMessageHandler()

	listenAgentChannel := false
	fmt.Println("started listening", a.id)

listening:
	for {
		select {
		case serverMessage := <-serverAgentchannel:
			//fmt.Println("server message", a.id, " ", serverMessage)
			switch serverMessage {
			case StartListeningNotification:
				//fmt.Println("started listening", a.id)
				listenAgentChannel = true
			case EndListeningNotification:
				//fmt.Println("stopped listening", a.id)
				listenAgentChannel = false
			case StopListeningSpinner:
				//fmt.Println("stopping listening on channel", a.id)
				break listening
			default:
				//fmt.Println("unknown message type")
			}
			a.AcknowledgeServerMessageReceived()
		default:
			if listenAgentChannel {
				select {
				case msg := <-agentAgentchannel:
					msg.Print()
					msg.InvokeMessageHandler(a.id)
				default:
				}
			}
		}
	}
	a.agentStoppedTalking(a.id)
	go a.AcknowledgeClosure(a.id)
	fmt.Println("stopped listening on channel")
}
