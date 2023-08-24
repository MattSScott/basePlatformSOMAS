package serverTesting

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/agent1"
	"github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/agent2"
	baseUserAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
	infra "github.com/MattSScott/basePlatformSOMAS/pkg/infra/server"
	testserver "github.com/MattSScott/basePlatformSOMAS/pkg/testServer"

	"testing"

	"github.com/google/uuid"
)

func TestBaseServer(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[baseAgent.Agent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(baseAgent.GetAgent, 4)

	serv := infra.CreateServer[baseAgent.Agent](m, 5)

	for _, agent := range serv.Agents {

		if agent.GetID() == uuid.Nil {
			t.Error("Error creating agent")

		}

	}
	serv.RunGameLoop()
	serv.Start()

}

func TestInheritedServer(t *testing.T) {

	m := make([]infra.AgentGeneratorCountPair[baseUserAgent.AgentUserInterface], 2)
	m[0] = infra.MakeAgentGeneratorCountPair[baseUserAgent.AgentUserInterface](agent2.GetAgent, 3)
	m[1] = infra.MakeAgentGeneratorCountPair[baseUserAgent.AgentUserInterface](agent1.GetAgent, 2)

	floors := 3
	ts := testserver.New(m, floors)

	if len(ts.Agents) != 5 {
		t.Error("Agents not properly instantiated")
	}

	if ts.GetName() != "TestServer" {
		t.Error("Server name not properly instantiated")
	}

	ts.RunGameLoop()
	ts.Start()
}
