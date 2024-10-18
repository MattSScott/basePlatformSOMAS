package message_test

import (
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/internal/testUtils"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

func TestMessageCanBeExtended(t *testing.T) {
	server := server.CreateBaseServer[testUtils.IExtendedAgent](1, 1, time.Second, 100000)
	agent1 := testUtils.NewExtendedAgent(server)
	agent2 := testUtils.NewExtendedAgent(server)
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
	server := server.CreateBaseServer[testUtils.IExtendedAgent](1, 1, time.Second, 100000)
	a1 := testUtils.NewExtendedAgent(server)
	a2 := testUtils.NewExtendedAgent(server)
	a3 := testUtils.NewExtendedAgent(server)
	nullMsg := a1.GetNullMessage()
	for _, recip := range []testUtils.IExtendedAgent{a2, a3} {
		nullMsg.InvokeMessageHandler(recip)
	}
	if a2.GetAgentField() != 0 || a3.GetAgentField() != 0 {
		t.Error("Message not properly distributed to recipients")
	}
}

func TestMultipleMessagesGetHandled(t *testing.T) {
	server := server.CreateBaseServer[testUtils.IExtendedAgent](1, 1, time.Second, 100000)
	a1 := testUtils.NewExtendedAgent(server)
	a2 := testUtils.NewExtendedAgent(server)
	allMessages := []message.IMessage[testUtils.IExtendedAgent]{a1.GetMessage1(), a1.GetMessage1(), a1.GetMessage2()}
	for _, msg := range allMessages {
		msg.InvokeMessageHandler(a2)
	}
	if a2.GetAgentField() != 20 {
		t.Error("Agents unable to handle multiple message types")
	}
}

func TestGetSender(t *testing.T) {
	server := server.CreateBaseServer[testUtils.IExtendedAgent](1, 1, time.Second, 100000)
	a1 := testUtils.NewExtendedAgent(server)
	msg := a1.GetMessage1()
	agID := a1.GetID()
	msgSenderID := msg.GetSender()
	if agID != msgSenderID {
		t.Errorf("Message has sender ID: %s, expected: %s", msgSenderID, agID)
	}
}
