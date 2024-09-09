package agent_test

import (
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/server"
	"github.com/google/uuid"
)

type IBaseAgent interface {
	agent.IAgent[IBaseAgent]
}

type AgentTestServer struct {
	*server.BaseServer[IBaseAgent]
}

func TestAgentIdOperations(t *testing.T) {
	var testServ server.IServer[IBaseAgent] = AgentTestServer{}
	baseAgent := agent.CreateBaseAgent[IBaseAgent](testServ)

	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
	}
}

type AgentWithState struct {
	*agent.BaseAgent[IBaseAgent]
	state int
}

func (aws *AgentWithState) UpdateAgentInternalState() {
	aws.state += 1
}

func TestUpdateAgentInternalState(t *testing.T) {
	var testServ server.IServer[IBaseAgent] = AgentTestServer{}

	ag := AgentWithState{
		agent.CreateBaseAgent[IBaseAgent](testServ),
		0,
	}

	if ag.state != 0 {
		t.Error("Additional agent field not correctly instantiated")
	}

	ag.UpdateAgentInternalState()

	if ag.state != 1 {
		t.Error("Agent state not correctly updated")
	}
}
