package baseserver_test

import (
	"testing"

	baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseAgent"
	baseserver "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseServer"
)

type ITestBaseAgent interface {
	baseagent.IAgent[ITestBaseAgent]
}

type TestBaseAgent struct {
	*baseagent.BaseAgent[ITestBaseAgent]
}

func NewTestBaseAgent() ITestBaseAgent {
	return &TestBaseAgent{
		baseagent.NewAgent[ITestBaseAgent](),
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	m := make([]baseserver.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = baseserver.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

	server := baseserver.CreateServer[ITestBaseAgent](m, 1)

	if len(server.GetAgents()) != 3 {
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
