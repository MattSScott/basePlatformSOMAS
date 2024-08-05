package infra_test

import (
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/infra"
)

type ITestBaseAgent interface {
	infra.IAgent[ITestBaseAgent]
}

type TestBaseAgent struct {
	*infra.BaseAgent[ITestBaseAgent]
}

func NewTestBaseAgent() ITestBaseAgent {
	var testServ infra.IServer[IBaseAgent] = TestServer{}

	return &TestBaseAgent{
		infra.CreateBaseAgent[IBaseAgent](testServ),
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	server := infra.CreateServer[ITestBaseAgent](m, 1)

	if len(server.GetAgentMap()) != 3 {
		t.Error("Incorrect number of agents added to server")
	}

}

func TestNumIterationsInServer(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	server := infra.CreateServer[ITestBaseAgent](m, 1)

	if server.GetIterations() != 1 {
		t.Error("Incorrect number of iterations instantiated")
	}

}

type IExtendedTestServer interface {
	infra.IServer[ITestBaseAgent]
	GetAdditionalField() int
}

type ExtendedTestServer struct {
	*infra.BaseServer[ITestBaseAgent]
	testField int
}

func (ets *ExtendedTestServer) GetAdditionalField() int {
	return ets.testField
}

func (ets *ExtendedTestServer) RunGameLoop() {

	ets.BaseServer.RunGameLoop()
	ets.testField += 1
}

func CreateTestServer(mapper []infra.AgentGeneratorCountPair[ITestBaseAgent], iters int) IExtendedTestServer {
	return &ExtendedTestServer{
		BaseServer: infra.CreateServer[ITestBaseAgent](mapper, iters),
		testField:  0,
	}
}

func TestAddAgent(t *testing.T) {

	baseServer := infra.CreateServer[ITestBaseAgent]([]infra.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)

	agent1 := infra.NewBaseAgent[ITestBaseAgent]()

	baseServer.AddAgent(agent1)

	if len(baseServer.GetAgentMap()) != 1 {
		t.Error("Agent not correctly added to map")
	}
}

func TestRemoveAgent(t *testing.T) {

	baseServer := infra.CreateServer[ITestBaseAgent]([]infra.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)

	agent1 := infra.NewBaseAgent[ITestBaseAgent]()

	baseServer.AddAgent(agent1)
	baseServer.RemoveAgent(agent1)

	if len(baseServer.GetAgentMap()) != 0 {
		t.Error("Agent not correctly removed from map")
	}
}

func TestFullAgentHashmap(t *testing.T) {
	baseServer := infra.CreateServer[ITestBaseAgent]([]infra.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)
	for i := 0; i < 5; i++ {
		baseServer.AddAgent(infra.NewBaseAgent[ITestBaseAgent]())
	}

	for id, agent := range baseServer.GetAgentMap() {
		if agent.GetID() != id {
			t.Error("Server agent hashmap key doesn't match object")
		}
	}
}

func TestServerGameLoop(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	server := CreateTestServer(m, 1)

	if server.GetAdditionalField() != 0 {
		t.Error("Additional server parameter not correctly instantiated")
	}

	server.RunGameLoop()

	if server.GetAdditionalField() != 1 {
		t.Error("Run Game Loop method not successfully overridden")
	}

}

func TestServerStartsCorrectly(t *testing.T) {
	generator := infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	baseServer := infra.CreateServer([]infra.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

	baseServer.Start()
}

func TestAgentMapConvertsToArray(t *testing.T) {
	generator := infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	baseServer := infra.CreateServer([]infra.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

	if len(baseServer.GenerateAgentArrayFromMap()) != 3 {
		t.Error("Agents not correctly mapped to array")
	}
}

func (tba *TestBaseAgent) GetAllMessages(others []ITestBaseAgent) []infra.IMessage[ITestBaseAgent] {
	msg := infra.CreateMessage[ITestBaseAgent](tba, others)

	return []infra.IMessage[ITestBaseAgent]{msg}
}

func TestMessagingSession(t *testing.T) {
	generator := infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	baseServer := infra.CreateServer([]infra.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

	agentArray := baseServer.GenerateAgentArrayFromMap()

	agent1 := agentArray[0]

	messages := agent1.GetAllMessages(agentArray)

	for _, msg := range messages {
		if len(msg.GetRecipients()) != 3 {
			t.Error("Incorrect number of message recipients")
		}
		for _, recip := range msg.GetRecipients() {
			if recip.GetID() == agent1.GetID() {
				continue
			}
			msg.InvokeMessageHandler(recip)
		}
	}

}
