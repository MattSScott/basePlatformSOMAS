package worldagent

import (
	"fmt"

	"github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
)

type WorldAgent struct {
	*baseExtendedAgent.BaseExtendedAgent
}

func (wa *WorldAgent) GetAllMessages(availableAgents []baseExtendedAgent.IExtendedAgent) []messaging.IMessage[baseExtendedAgent.IExtendedAgent] {

	msg := wa.CreateGreetingMessage(availableAgents)

	return []messaging.IMessage[baseExtendedAgent.IExtendedAgent]{msg}

}

func (wa *WorldAgent) CreateGreetingMessage(recips []baseExtendedAgent.IExtendedAgent) baseExtendedAgent.GreetingMessage {
	return baseExtendedAgent.CreateGreetingMessage(wa, recips, wa.GetPhrase())
}

func (wa *WorldAgent) HandleGreetingMessage(msg baseExtendedAgent.GreetingMessage) {
	if msg.GetGreeting() == "hello" {
		fmt.Println(wa.GetID(), " responded: 'world'")
	}
}

func GetWorldAgent() baseExtendedAgent.IExtendedAgent {
	return &WorldAgent{
		BaseExtendedAgent: baseExtendedAgent.GetAgent("wello"),
	}
}
