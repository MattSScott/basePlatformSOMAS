package message_test

import (
	"testing"

	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
)

type Message1 struct {
	message.BaseMessage[IExtendedAgent]
	messageField1 int
}

type Message2 struct {
	message.BaseMessage[IExtendedAgent]
	messageField2 int
}

type NullMessage struct {
	message.BaseMessage[IExtendedAgent]
}

type IExtendedAgent interface {
	message.IAgentMessenger[IExtendedAgent]
	GetAgentField() int
	GetAllMessages() []message.IMessage[IExtendedAgent]
	HandleMessage1(msg Message1)
	HandleMessage2(msg Message2)
	HandleNullMessage(msg NullMessage)
}

type ExtendedAgent struct {
	agentField int
}

func (m1 Message1) Accept(agent IExtendedAgent) {
	agent.HandleMessage1(m1)
}

func (m2 Message2) Accept(agent IExtendedAgent) {
	agent.HandleMessage2(m2)
}

func (nm NullMessage) Accept(agent IExtendedAgent) {
	agent.HandleNullMessage(nm)
}

func (ea *ExtendedAgent) GetAgentField() int {
	return ea.agentField
}

func (ea *ExtendedAgent) GetMessage1() Message1 {
	return Message1{
		message.BaseMessage[IExtendedAgent]{},
		5,
	}
}

func (ea *ExtendedAgent) HandleMessage1(msg Message1) {
	ea.agentField += msg.messageField1
}

func (ea *ExtendedAgent) GetMessage2() Message2 {
	return Message2{
		message.BaseMessage[IExtendedAgent]{},
		10,
	}
}

func (ea *ExtendedAgent) HandleMessage2(msg Message2) {
	ea.agentField += msg.messageField2
}
func (ea *ExtendedAgent) GetNullMessage(recips []IExtendedAgent) NullMessage {
	return NullMessage{message.CreateMessage[IExtendedAgent](ea, recips)}
}

func (ea *ExtendedAgent) HandleNullMessage(msg NullMessage) {
	sender := msg.GetSender()
	ea.agentField = sender.GetAgentField()
}

func (ea *ExtendedAgent) GetAllMessages() []message.IMessage[IExtendedAgent] {
	msg1 := ea.GetMessage1()
	msg2 := ea.GetMessage2()
	return []message.IMessage[IExtendedAgent]{msg1, msg2}
}
func TestMessageCanBeExtended(t *testing.T) {
	agent1 := &ExtendedAgent{agentField: 0}
	agent2 := &ExtendedAgent{agentField: 0}

	msgFromA1 := agent1.GetMessage1()
	msgFromA1.Accept(agent2)

	msgFromA2 := agent2.GetMessage2()
	msgFromA2.Accept(agent1)

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
		nullMsg.Accept(recip)
	}

	if a1.GetAgentField() != a2.GetAgentField() || a3.GetAgentField() != a1.GetAgentField() {
		t.Error("Message not properly distributed to recipients")
	}

}
