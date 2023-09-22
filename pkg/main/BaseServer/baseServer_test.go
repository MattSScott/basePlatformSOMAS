package baseserver_test

import (
	"testing"

	baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseAgent"
	baseserver "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseServer"
	"github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"
)

type ITestBaseAgent interface {
	baseagent.IAgent[ITestBaseAgent]
}

type TestBaseAgent struct {
	*baseagent.BaseAgent[ITestBaseAgent]
}

func NewTestBaseAgent() ITestBaseAgent {
	return &TestBaseAgent{
		baseagent.NewBaseAgent[ITestBaseAgent](),
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	m := make([]baseserver.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = baseserver.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	server := baseserver.CreateServer[ITestBaseAgent](m, 1)

	if len(server.GetAgentMap()) != 3 {
		t.Error("Incorrect number of agents added to server")
	}

}

func TestNumTurnsInServer(t *testing.T) {
	m := make([]baseserver.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = baseserver.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	server := baseserver.CreateServer[ITestBaseAgent](m, 1)

	if server.GetNumTurns() != 1 {
		t.Error("Incorrect number of turns instantiated")
	}

}

type IExtendedTestServer interface {
	baseserver.IServer[ITestBaseAgent]
	GetAdditionalField() int
}

type ExtendedTestServer struct {
	*baseserver.BaseServer[ITestBaseAgent]
	testField int
}

func (ets *ExtendedTestServer) GetAdditionalField() int {
	return ets.testField
}

func (ets *ExtendedTestServer) RunGameLoop() {

	ets.BaseServer.RunGameLoop()
	ets.testField += 1
}

func CreateTestServer(mapper []baseserver.AgentGeneratorCountPair[ITestBaseAgent], iters int) IExtendedTestServer {
	return &ExtendedTestServer{
		BaseServer: baseserver.CreateServer[ITestBaseAgent](mapper, iters),
		testField:  0,
	}
}

func TestAddAgent(t *testing.T) {

	baseServer := baseserver.CreateServer[ITestBaseAgent]([]baseserver.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)

	agent1 := baseagent.NewBaseAgent[ITestBaseAgent]()

	baseServer.AddAgent(agent1)

	if len(baseServer.GetAgentMap()) != 1 {
		t.Error("Agent not correctly added to map")
	}
}

func TestRemoveAgent(t *testing.T) {

	baseServer := baseserver.CreateServer[ITestBaseAgent]([]baseserver.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)

	agent1 := baseagent.NewBaseAgent[ITestBaseAgent]()

	baseServer.AddAgent(agent1)
	baseServer.RemoveAgent(agent1)

	if len(baseServer.GetAgentMap()) != 0 {
		t.Error("Agent not correctly removed from map")
	}
}

func TestFullAgentHashmap(t *testing.T) {
	baseServer := baseserver.CreateServer[ITestBaseAgent]([]baseserver.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)
	for i := 0; i < 5; i++ {
		baseServer.AddAgent(baseagent.NewBaseAgent[ITestBaseAgent]())
	}

	for id, agent := range baseServer.GetAgentMap() {
		if agent.GetID() != id {
			t.Error("Server agent hashmap key doesn't match object")
		}
	}
}

func TestServerGameLoop(t *testing.T) {
	m := make([]baseserver.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = baseserver.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

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
	generator := baseserver.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	baseServer := baseserver.CreateServer([]baseserver.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

	baseServer.Start()
}

func TestAgentMapConvertsToArray(t *testing.T) {
	generator := baseserver.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	baseServer := baseserver.CreateServer([]baseserver.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

	if len(baseServer.GenerateAgentArrayFromMap()) != 3 {
		t.Error("Agents not correctly mapped to array")
	}
}

func (tba *TestBaseAgent) GetAllMessages(others []ITestBaseAgent) []messaging.IMessage[ITestBaseAgent] {
	msg := messaging.CreateMessage[ITestBaseAgent](tba, others)

	return []messaging.IMessage[ITestBaseAgent]{msg}
}

func TestMessagingSession(t *testing.T) {
	generator := baseserver.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	baseServer := baseserver.CreateServer([]baseserver.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

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
