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
	msgText := fmt.Sprintf("%s said: '%s'", ha.GetID(), ha.GetPhrase())
	fmt.Println(msgText)
	return baseExtendedAgent.CreateGreetingMessage(ha, recips, ha.GetPhrase())
}

func (ha *HelloAgent) HandleGreetingMessage(msg baseExtendedAgent.GreetingMessage) {
	if msg.GetGreeting() == "wello" {
		respText := fmt.Sprintf("%s responded: 'horld'", ha.GetID())
		fmt.Println(respText)
	}
}

func GetHelloAgent() baseExtendedAgent.IExtendedAgent {
	return &HelloAgent{
		BaseExtendedAgent: baseExtendedAgent.GetBaseExtendedAgent("hello"),
	}
}
