package message_test

import (
	"fmt"
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/pkg/internal/testUtils"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
)

func TestMessageCanBeExtended(t *testing.T) {
	agent1 := testUtils.NewExtendedAgent(nil)
	agent2 := testUtils.NewExtendedAgent(nil)

	msgFromA1 := agent1.GetMessage1()
	msgFromA1.InvokeMessageHandler(agent2)

	msgFromA2 := agent2.GetMessage2()
	msgFromA2.InvokeMessageHandler(agent1)

	if agent1.GetAgentField() != 10 {
		t.Error("Message 2 not properly handled")
	}

	if agent2.GetAgentField() != 5 {
		t.Error("Message 1 not properly handled")
	}
}

func TestSingleMessageGetsHandled(t *testing.T) {

	a1 := testUtils.NewExtendedAgent(nil)
	a2 := testUtils.NewExtendedAgent(nil)
	a3 := testUtils.NewExtendedAgent(nil)

	nullMsg := a1.GetNullMessage()

	for _, recip := range []testUtils.IExtendedAgent{a2, a3} {
		nullMsg.InvokeMessageHandler(recip)
	}

	//fmt.Println(a2.AgentField)

	if a2.GetAgentField() != 0 || a3.GetAgentField() != 0 {
		t.Error("Message not properly distributed to recipients")
	}
}

func TestMultipleMessagesGetHandled(t *testing.T) {

	a1 := testUtils.NewExtendedAgent(nil)
	a2 := testUtils.NewExtendedAgent(nil)

	allMessages := []message.IMessage[testUtils.IExtendedAgent]{a1.GetMessage1(), a1.GetMessage1(), a1.GetMessage2()}

	for _, msg := range allMessages {
		msg.InvokeMessageHandler(a2)
	}

	fmt.Println(a2.GetAgentField())

	if a2.GetAgentField() != 20 {
		t.Error("Agents unable to handle multiple message types")
	}
}

func TestPrint(t *testing.T) {
	a1 := testUtils.NewExtendedAgent(nil)
	newMsg := a1.GetMessage1()
	newMsg.Print()
}

func TestGetSender(t *testing.T) {
	a1 := testUtils.NewExtendedAgent(nil)
	msg:=a1.GetMessage1()
	agID := a1.GetID()
	msgSenderID := msg.GetSender()
	if agID != msgSenderID {
		t.Errorf("Message has sender ID: %s, expected: %s",msgSenderID, agID)
	}
}
