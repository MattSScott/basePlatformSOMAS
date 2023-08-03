package agentTesting

import (
	agent1 "basePlatformSOMAS/pkg/agents/AgentTesting/agent1"
	agent2 "basePlatformSOMAS/pkg/agents/AgentTesting/agent2"
	"testing"
)

func TestAgentOneName(t *testing.T) {
	a1 := agent1.GetAgent().(*agent1.Agent1)
	if a1.GetName() != "A1" {
		t.Error("Expected name to be A1")
	}

}

func TestAgentOneAge(t *testing.T) {
	a1 := agent1.GetAgent().(*agent1.Agent1)

	beforeAge := a1.GetAge()
	a1.UpdateAgent()
	afterAge := a1.GetAge()

	if afterAge != beforeAge+1 {
		t.Error("Age not increased after update")
	}

}

func TestAgentTwo(t *testing.T) {
	a2 := agent2.GetAgent().(*agent2.Agent2)

	gender := a2.GetGender()

	if gender != "Male" && gender != "Female" {
		t.Error("Gender not appropriately set")
	}

}
