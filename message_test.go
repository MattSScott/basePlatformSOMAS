package basePlatformSOMAS_test

import (
	"fmt"
	"testing"

	"github.com/MattSScott/basePlatformSOMAS"
)

type Message1 struct {
	basePlatformSOMAS.BaseMessage
	messageField1 int
}

type Message2 struct {
	basePlatformSOMAS.BaseMessage
	messageField2 int
}

type NullMessage struct {
	basePlatformSOMAS.BaseMessage
}

type IExtendedAgent interface {
	basePlatformSOMAS.IAgent[IExtendedAgent]
	GetAgentField() int
	HandleMessage1(msg Message1)
	HandleMessage2(msg Message2)
	HandleNullMessage(msg NullMessage)
}

type ExtendedAgent struct {
	basePlatformSOMAS.BaseAgent[IExtendedAgent]
	agentField int
}

func (agent *ExtendedAgent) HandleMessage1(msg Message1) {
	agent.agentField += msg.messageField1
}

func (agent *ExtendedAgent) HandleMessage2(msg Message2) {
	agent.agentField += msg.messageField2

}

func (agent *ExtendedAgent) HandleNullMessage(msg NullMessage) {
	fmt.Println("handled!")
	agent.agentField = 0
	fmt.Println(agent.agentField)
}

func (m1 Message1) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage1(m1)
}

func (m2 Message2) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage2(m2)
}

func (nm NullMessage) InvokeMessageHandler(agent IExtendedAgent) {
	fmt.Println("Handling!")
	agent.HandleNullMessage(nm)
}

func (m1 Message1) InvokeSyncMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage1(m1)
}

func (m2 Message2) InvokeSyncMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage2(m2)
}

func (nm NullMessage) InvokeSyncMessageHandler(agent IExtendedAgent) {
	fmt.Println("Handling!")
	agent.HandleNullMessage(nm)
}

func (ea *ExtendedAgent) GetAgentField() int {
	return ea.agentField
}

func (ea *ExtendedAgent) GetMessage1() Message1 {
	return Message1{
		ea.CreateBaseMessage(),
		5,
	}
}

func (ea *ExtendedAgent) GetMessage2() Message2 {
	return Message2{
		ea.CreateBaseMessage(),
		10,
	}
}

func (ea *ExtendedAgent) GetNullMessage() NullMessage {
	return NullMessage{
		ea.CreateBaseMessage(),
	}
}

func TestMessageCanBeExtended(t *testing.T) {
	agent1 := &ExtendedAgent{agentField: 0}
	agent2 := &ExtendedAgent{agentField: 0}

	msgFromA1 := agent1.GetMessage1()
	msgFromA1.InvokeMessageHandler(agent2)

	msgFromA2 := agent2.GetMessage2()
	msgFromA2.InvokeMessageHandler(agent1)

	if agent1.agentField != 10 {
		t.Error("Message 2 not properly handled")
	}

	if agent2.agentField != 5 {
		t.Error("Message 1 not properly handled")
	}
}
func TestSingleMessageGetsHandled(t *testing.T) {

	a1 := &ExtendedAgent{agentField: 5}
	a2 := &ExtendedAgent{agentField: 10}
	a3 := &ExtendedAgent{agentField: 15}

	nullMsg := a1.GetNullMessage()

	for _, recip := range []IExtendedAgent{a2, a3} {
		nullMsg.InvokeMessageHandler(recip)
	}

	fmt.Println(a2.agentField)

	if a2.GetAgentField() != 0 || a3.GetAgentField() != 0 {
		t.Error("Message not properly distributed to recipients")
	}
}

func TestMultipleMessagesGetHandled(t *testing.T) {

	a1 := &ExtendedAgent{agentField: 0}
	a2 := &ExtendedAgent{agentField: 0}

	allMessages := []basePlatformSOMAS.IMessage[IExtendedAgent]{a1.GetMessage1(), a1.GetMessage1(), a1.GetMessage2()}

	for _, msg := range allMessages {
		msg.InvokeMessageHandler(a2)
	}

	fmt.Println(a2.agentField)

	if a2.GetAgentField() != 20 {
		t.Error("Agents unable to handle multiple message types")
	}
}
