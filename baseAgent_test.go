package basePlatformSOMAS_test

import (
	"testing"

	"github.com/MattSScott/basePlatformSOMAS"
	"github.com/google/uuid"
)

type IBaseAgent interface {
	basePlatformSOMAS.IAgent[IBaseAgent]
}

type AgentTestServer struct {
	*basePlatformSOMAS.BaseServer[IBaseAgent]
}

func TestAgentIdOperations(t *testing.T) {
	var testServ basePlatformSOMAS.IServer[IBaseAgent] = AgentTestServer{}
	baseAgent := basePlatformSOMAS.CreateBaseAgent[IBaseAgent](testServ)

	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
	}
}

type AgentWithState struct {
	*basePlatformSOMAS.BaseAgent[IBaseAgent]
	state int
}

func (aws *AgentWithState) UpdateAgentInternalState() {
	aws.state += 1
}

func TestUpdateAgentInternalState(t *testing.T) {
	var testServ basePlatformSOMAS.IServer[IBaseAgent] = AgentTestServer{}

	ag := AgentWithState{
		basePlatformSOMAS.CreateBaseAgent[IBaseAgent](testServ),
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
