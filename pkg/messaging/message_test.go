package message_test

import (
	"fmt"
	"testing"

	baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
)

type ComplexMessage struct {
	message.BaseMessage
	additionalInformation string
}

type IComplexAgent interface {
	baseagent.IAgent
	// message.IAgentMessaging

	HandleComplexMessage(msg ComplexMessage)
}

type ComplexAgent struct {
	baseagent.BaseAgent
}

func (cm ComplexMessage) GetAdditionalInformation() string {
	return cm.additionalInformation
}

func (ca *ComplexAgent) GetMessage() message.IMessage[IComplexAgent] {
	return CreateComplexMessage(ca, "additionalInfo")
}

// func (ca *ComplexAgent) HandleMessage(msg message.IMessage) {
// 	ca.HandleComplexMessage(msg)
// }

func (ca *ComplexAgent) HandleComplexMessage(msg ComplexMessage) {
	fmt.Println("Complex message received")
}

func (cm ComplexMessage) HowToHandleMessage(agent IComplexAgent) {

	// var ca IComplexAgent = agent

	agent.HandleComplexMessage(cm)

}

func CreateComplexMessage(sender IComplexAgent, addedInfo string) ComplexMessage {
	return ComplexMessage{
		message.CreateNullMessageWithSender(sender),
		addedInfo,
	}
}

func TestMessageComposition(t *testing.T) {
	// ca := ComplexAgent{}
	// msg := ca.GetMessage()
	// content := msg.GetAdditionalInformation()

}
