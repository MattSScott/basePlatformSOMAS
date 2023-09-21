package baseagent_test

import (
	"testing"

	baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseAgent"
	"github.com/google/uuid"
)

type IBaseAgent interface {
	baseagent.IAgent[IBaseAgent]
}

func TestAgentIdOperations(t *testing.T) {
	baseAgent := baseagent.NewAgent[IBaseAgent]()

	if baseAgent.GetID() == uuid.Nil {
		t.Error("Agent not instantiated with valid ID")
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
