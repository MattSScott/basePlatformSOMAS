package agent_test

import (
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/internal/testUtils"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/google/uuid"
)

func TestAgentIdOperations(t *testing.T) {
	var testServ agent.IExposedServerFunctions[testUtils.ITestBaseAgent] = &testUtils.TestServer{
		BaseServer:            server.CreateBaseServer[testUtils.ITestBaseAgent](1, 1, time.Second, 100),
		TurnCounter:           0,
		IterationStartCounter: 0,
		IterationEndCounter:   0,
	}
	baseAgent := agent.CreateBaseAgent(testServ)
	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
	}
}

func TestNilInterfaceInjection(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("did not panic when nil interface injected")
		}
	}()
	ag := agent.CreateBaseAgent[testUtils.ITestBaseAgent](nil)
	ag.GetID()
}

func TestCreateBaseMessage(t *testing.T) {
	testServ := testUtils.GenerateTestServer(1, 1, 1, time.Second, 100000)
	ag := testUtils.NewTestAgent(testServ)
	newMsg := ag.CreateBaseMessage()
	msgSenderID := newMsg.GetSender()
	agID := ag.GetID()
	if msgSenderID != agID {
		t.Error("Incorrect Sender ID in Message. Expected:", agID, "got:", msgSenderID)
	}
}

func TestNotifyAgentMessaging(t *testing.T) {
	testServ := testUtils.GenerateTestServer(1, 1, 1, time.Second, 100000)
	ag := testUtils.NewTestAgent(testServ)
	ag.FinishedMessaging()
	agentStoppedTalkingCalls := ag.GetAgentStoppedTalking()
	if agentStoppedTalkingCalls != 1 {
		t.Error("expected 1 calls of agentStoppedTalking(), got:", agentStoppedTalkingCalls)
	}
}

func TestBroadcastMessageNoIDPanic(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("did not panic when message sender not set")
		}
	}()
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Millisecond, 100000)
	for _, ag := range server.GetAgentMap() {
		msg := &testUtils.TestMessage{}
		ag.BroadcastMessage(msg)
	}
}

func TestSendMessageNoIDPanic(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("did not panic when message sender not set")
		}
	}()
	numAgents := 1
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Millisecond, 100000)
	agMap := server.GetAgentMap()

	for _, ag := range agMap {
		msg := &testUtils.TestMessage{}
		for recip := range agMap {
			ag.SendMessage(msg, recip)
		}
	}
}

func TestSendSynchronousMessageNoIDPanic(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("did not panic when message sender not set")
		}
	}()
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Millisecond, 100)
	agMap := server.GetAgentMap()
	for _, ag := range agMap {
		msg := &testUtils.TestMessage{}
		for recip := range agMap {
			ag.SendSynchronousMessage(msg, recip)
		}
	}
}

func TestSendSynchronousMessage(t *testing.T) {
	numAgents := 10
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	testMessage := testUtils.NewTestAgent(server).CreateTestMessage()
	for id, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
		ag.SendSynchronousMessage(testMessage, id)
	}
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error("Didn't Receive Message")
		}
	}
}

func TestRecursiveInvokeMessageHandlerCalls(t *testing.T) {
	numAgents := 3
	timeLimit := time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	for _, ag := range server.GetAgentMap() {
		msg := testUtils.CreateInfLoopMessage(ag.GetID())
		ag.BroadcastMessage(msg)
	}
	time.Sleep(1 * time.Second)
}

func TestSendMessage(t *testing.T) {
	numAgents := 3
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100000)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	server.AddAgent(agent1)
	for id, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
		agent1.SendMessage(testMessage, id)
	}
	time.Sleep(1 * time.Second)
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error(ag, "Didn't Receive Message")
		}
	}
}

func TestBroadcastMessage(t *testing.T) {
	numAgents := 3
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100000)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	server.AddAgent(agent1)
	for _, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
	}
	agent1.BroadcastMessage(testMessage)
	senderID := agent1.GetID()
	time.Sleep(1 * time.Second)
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() && ag.GetID() != senderID {
			t.Error(ag, "Didn't Receive Message")
		} else if ag.ReceivedMessage() && ag.GetID() == senderID {
			t.Error(ag, "is sender and received its own message")
		}
	}
}

func TestBroadcastMessageDoesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("BroadcastMessage did not panic as expected")
		}
	}()
	numAgents := 3
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := &testUtils.TestMessage{}
	server.AddAgent(agent1)
	for _, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
	}
	agent1.BroadcastMessage(testMessage)
	time.Sleep(1 * time.Second)
}

func TestBroadcastSynchronousMessage(t *testing.T) {
	numAgents := 3
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	server.AddAgent(agent1)
	for _, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
	}
	agent1.BroadcastSynchronousMessage(testMessage)
	senderID := agent1.GetID()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() && ag.GetID() != senderID {
			t.Error(ag, "Didn't Receive Message")
		} else if ag.ReceivedMessage() && ag.GetID() == senderID {
			t.Error(ag, "is sender and received its own message")
		}
	}
}

func TestBroadcastSynchronousMessageDoesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("BroadcastSynchronousMessage did not panic as expected")
		}
	}()
	numAgents := 3
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := &testUtils.TestMessage{}
	server.AddAgent(agent1)
	for _, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
	}
	agent1.BroadcastSynchronousMessage(testMessage)
}
