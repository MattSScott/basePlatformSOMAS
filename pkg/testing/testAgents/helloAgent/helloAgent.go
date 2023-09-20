package helloagent

import (
	"fmt"

	"github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
)

type HelloAgent struct {
	*baseExtendedAgent.BaseExtendedAgent
}

func (ha *HelloAgent) GetAllMessages(availableAgents []baseExtendedAgent.IExtendedAgent) []messaging.IMessage[baseExtendedAgent.IExtendedAgent] {

	msg := ha.CreateGreetingMessage(availableAgents)

	return []messaging.IMessage[baseExtendedAgent.IExtendedAgent]{msg}

}

func (ha *HelloAgent) CreateGreetingMessage(recips []baseExtendedAgent.IExtendedAgent) baseExtendedAgent.GreetingMessage {
	return baseExtendedAgent.CreateGreetingMessage(ha, recips, ha.GetPhrase())
}

func (ha *HelloAgent) HandleGreetingMessage(msg baseExtendedAgent.GreetingMessage) {
	if msg.GetGreeting() == "wello" {
		fmt.Println(ha.GetID(), " responded: 'horld'")
	}
}

func GetHelloAgent() baseExtendedAgent.IExtendedAgent {
	return &HelloAgent{
		BaseExtendedAgent: baseExtendedAgent.GetAgent("hello"),
	}
}
