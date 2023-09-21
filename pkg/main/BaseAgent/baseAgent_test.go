package baseagent_test

import (
	"testing"

	baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseAgent"
	"github.com/google/uuid"
)

type IBaseAgent interface {
	baseagent.IAgent[IBaseAgent]
}

func TestAgentIDOperations(t *testing.T) {
	baseAgent := baseagent.NewAgent[IBaseAgent]()

	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
	}
}

func TestAddToNetwork(t *testing.T) {
	agent1 := baseagent.NewAgent[IBaseAgent]()
	agent2 := baseagent.NewAgent[IBaseAgent]()

	agent1.AddAgentToNetwork(agent2)

	if len(agent1.GetNetwork()) != 1 {
		t.Error("Agent not correctly added to map")
	}
}

func TestRemoveFromNetwork(t *testing.T) {
	agent1 := baseagent.NewAgent[IBaseAgent]()
	agent2 := baseagent.NewAgent[IBaseAgent]()

	agent1.AddAgentToNetwork(agent2)
	agent1.RemoveAgentFromNetwork(agent2)

	if len(agent1.GetNetwork()) != 0 {
		t.Error("Agent not correctly removed from map")
	}
}

type AgentWithState struct {
	*baseagent.BaseAgent[*AgentWithState]
	state int
}

func (aws *AgentWithState) UpdateAgentInternalState() {
	aws.state += 1
}

func TestUpdateAgentInternalState(t *testing.T) {
	ag := AgentWithState{
		baseagent.NewAgent[*AgentWithState](),
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

	agent := baseagent.NewAgent[IBaseAgent]()

	messages := agent.GetAllMessages([]IBaseAgent{agent})

	if len(messages) > 0 {
		t.Error("Agent erroneously constructed message")
	}

}
