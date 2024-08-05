package infra_test

import (
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/infra"
	"github.com/google/uuid"
)

type IBaseAgent interface {
	infra.IAgent[IBaseAgent]
}

func TestAgentIdOperations(t *testing.T) {
	var IServ infra.IServer[IBaseAgent] = infra.IServer[IBaseAgent]{}
	baseAgent := infra.CreateBaseAgent(IServ)

	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
	}
}

type AgentWithState struct {
	*BaseAgent[*AgentWithState]
	state int
}

func (aws *AgentWithState) UpdateAgentInternalState() {
	aws.state += 1
}

func TestUpdateAgentInternalState(t *testing.T) {
	ag := AgentWithState{
		infra.CreateBaseAgent[*AgentWithState](),
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

func TestMessageRetrieval(t *testing.T) {

	agent := infra.CreateBaseAgent[IBaseAgent]()

	messages := agent.GetAllMessages([]IBaseAgent{agent})

	if len(messages) > 0 {
		t.Error("Agent erroneously constructed message")
	}

}
