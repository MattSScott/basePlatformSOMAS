package testserver_test

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/main/server"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
	helloagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/helloAgent"
	worldagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/worldAgent"
	testserver "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testServer"

	"testing"
)

func TestInheritedServer(t *testing.T) {

	m := make([]server.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], 5)

	m[0] = server.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](helloagent.GetHelloAgent, 3)
	m[1] = server.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](worldagent.GetWorldAgent, 2)

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
