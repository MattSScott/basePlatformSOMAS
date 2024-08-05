package infra

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type BaseAgent struct {
	unexportedServerFunctions
	IMessagingProtocol
	ListenOnChannel
	id uuid.UUID
}

func (ba *BaseAgent) GetID() uuid.UUID {
	return ba.id
}

func CreateBaseAgent() *BaseAgent {
	return &BaseAgent{
		id: uuid.New(),
	}
}

func (a *BaseAgent) NotifyAgentInactive() {
	a.agentStoppedTalking(a.id)
}

func (a *BaseAgent) listenOnChannel(agentAgentchannel chan IMessage, serverAgentchannel chan ServerNotification, wait *sync.WaitGroup) {
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
	go a.AcknowledgeClosure(a.id)
	fmt.Println("stopped listening on channel")
}
