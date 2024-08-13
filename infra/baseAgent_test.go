package infra_test

import (
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/infra"
	"github.com/google/uuid"
)

type IBaseAgent interface {
	infra.IAgent[IBaseAgent]
}

type AgentTestServer struct {
	*infra.BaseServer[IBaseAgent]
}

func TestAgentIdOperations(t *testing.T) {
	var testServ infra.IServer[IBaseAgent] = AgentTestServer{}
	baseAgent := infra.CreateBaseAgent[IBaseAgent](testServ)

	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
	}
}

type AgentWithState struct {
	*infra.BaseAgent[IBaseAgent]
	state int
}

func (aws *AgentWithState) UpdateAgentInternalState() {
	aws.state += 1
}

func TestUpdateAgentInternalState(t *testing.T) {
	var testServ infra.IServer[IBaseAgent] = AgentTestServer{}

	ag := AgentWithState{
		infra.CreateBaseAgent[IBaseAgent](testServ),
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
