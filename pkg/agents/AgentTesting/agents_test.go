package agentTesting

import (
	agent1 "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/agent1"
	agent2 "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/agent2"
	"testing"
)

func TestAgentOneName(t *testing.T) {
	a1 := agent1.GetAgent().(*agent1.Agent1)
	if a1.GetName() != "A1" {
		t.Error("Expected name to be A1")
	}

}

func TestAgentTwoName(t *testing.T) {
	a2 := agent2.GetAgent().(*agent2.Agent2)
	if a2.GetName() != "A2" {
		t.Error("Expected name to be A2")
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

func TestAgentOneFood(t *testing.T) {
	a1 := agent1.GetAgent().(*agent1.Agent1)

	foodInitial := a1.GetFood()

	if foodInitial != 0 {
		t.Error("Food not appropriately set")

	}

	a1.SetFood()
	food := a1.GetFood()
	if food != 1 {
		t.Error("Food not appropriately added")

	}
}

func TestAgentTwoFood(t *testing.T) {
	a2 := agent2.GetAgent().(*agent2.Agent2)

	foodInitial := a2.GetFood()

	if foodInitial != 0 {
		t.Error("Food not appropriately set")

	}

	a2.SetFood()
	food := a2.GetFood()
	if food != 1 {
		t.Error("Food not appropriately added")

	}
}

func TestAgentTwoSleep(t *testing.T) {
	a2 := agent2.GetAgent().(*agent2.Agent2)

	sleep := a2.GetSleep()

	if sleep != 100 {
		t.Error("Sleep not appropriately set")
	}

	a2.Activity()
	sleep = a2.GetSleep()

	if sleep != 90 {
		t.Error("Sleep not appropriately set")
	}

}
