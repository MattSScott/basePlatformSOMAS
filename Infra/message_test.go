package infra_test

import (
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/infra"
)

type Message1 struct {
	infra.BaseMessage
	messageField1 int
}

type Message2 struct {
	infra.BaseMessage
	messageField2 int
}

type NullMessage struct {
	infra.BaseMessage
}

type IExtendedAgent interface {
	infra.IAgent[IExtendedAgent]
	GetAgentField() int
	GetAllMessages([]IExtendedAgent) []infra.IMessage
	HandleMessage1(msg Message1)
	HandleMessage2(msg Message2)
	HandleNullMessage(msg NullMessage)
}

type ExtendedAgent struct {
	infra.BaseAgent[IExtendedAgent]
	agentField int
}

func (m1 Message1) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage1(m1)
}

func (m2 Message2) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage2(m2)
}

func (nm NullMessage) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleNullMessage(nm)
}

func (ea *ExtendedAgent) GetAgentField() int {
	return ea.agentField
}

func (ea *ExtendedAgent) GetMessage1() Message1 {
	return Message1{
		infra.BaseMessage{},
		5,
	}
}

func (ea *ExtendedAgent) HandleMessage1(msg Message1) {
	ea.agentField += msg.messageField1
}

func (ea *ExtendedAgent) GetMessage2() Message2 {
	return Message2{
		infra.BaseMessage{},
		10,
	}
}

func (ea *ExtendedAgent) HandleMessage2(msg Message2) {
	ea.agentField += msg.messageField2
}
func (ea *ExtendedAgent) GetNullMessage(recips []IExtendedAgent) NullMessage {
	return NullMessage{infra.CreateMessage[IExtendedAgent](ea, recips)}
}

func (ea *ExtendedAgent) HandleNullMessage(msg NullMessage) {
	// sender := msg.GetSender()

}

func (ea *ExtendedAgent) GetAllMessages([]IExtendedAgent) []infra.IMessage[IExtendedAgent] {
	msg1 := ea.GetMessage1()
	msg2 := ea.GetMessage2()
	return []infra.IMessage{msg1, msg2}
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

func TestMessageSender(t *testing.T) {

	a1 := &ExtendedAgent{agentField: 5}
	a2 := &ExtendedAgent{agentField: 10}
	a3 := &ExtendedAgent{agentField: 15}

	nullMsg := a1.GetNullMessage([]IExtendedAgent{a2, a3})

	for _, recip := range nullMsg.GetRecipients() {
		nullMsg.InvokeMessageHandler(recip)
	}

	if a1.GetAgentField() != a2.GetAgentField() || a3.GetAgentField() != a1.GetAgentField() {
		t.Error("Message not properly distributed to recipients")
	}

}
