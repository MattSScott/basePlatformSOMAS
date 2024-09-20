package agent_test

import (
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/internal/testUtils"
	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/google/uuid"
)

func TestAgentIdOperations(t *testing.T) {
	var testServ agent.IExposedServerFunctions[testUtils.ITestBaseAgent] = &testUtils.TestServer{}
	baseAgent := agent.CreateBaseAgent(testServ)
	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
	}
}

func TestUpdateAgentInternalState(t *testing.T) {
	var testServ agent.IExposedServerFunctions[testUtils.ITestBaseAgent] = &testUtils.TestServer{}
	ag := testUtils.TestServerFunctionsAgent{
		BaseAgent: agent.CreateBaseAgent(testServ),
		Counter:   0,
	}
	if ag.Counter != 0 {
		t.Error("Additional agent field not correctly instantiated")
	}
	ag.UpdateAgentInternalState()
	if ag.Counter != 1 {
		t.Error("Agent state not correctly updated")
	}
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
