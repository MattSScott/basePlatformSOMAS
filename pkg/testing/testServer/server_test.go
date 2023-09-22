package testserver_test

import (
	baseserver "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseServer"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
	helloagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/helloAgent"
	worldagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/worldAgent"
	testserver "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testServer"
	"github.com/google/uuid"

	"testing"
)

func TestInheritedServer(t *testing.T) {

	m := make([]baseserver.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], 2)

	m[0] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](helloagent.GetHelloAgent, 3)
	m[1] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](worldagent.GetWorldAgent, 2)

	floors := 3
	ts := testserver.New(m, floors)

	if len(ts.GetAgents()) != 5 {
		t.Error("Agents not properly instantiated")
	}

	if ts.GetName() != "TestServer" {
		t.Error("Server name not properly instantiated")
	}

	ts.RunGameLoop()
	ts.Start()
}

func TestServerInterfaceComposition(t *testing.T) {

	m := make([]baseserver.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], 2)

	m[0] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](helloagent.GetHelloAgent, 1)
	m[1] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](worldagent.GetWorldAgent, 1)

	// server can also be declared as type interface
	var server testserver.IExtendedServer = testserver.New(m, 1)

	if len(server.GetAgents()) != 2 {
		t.Error("Agents not properly instantiated")
	}

	for _, agent := range server.GetAgents() {
		if agent.GetID() == uuid.Nil {
			t.Error("Agent types not correctly instantiated")
		}
	}

	server.Start()
}

func TestServerMessagePassing(t *testing.T) {
	m := make([]baseserver.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], 2)

	m[0] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](helloagent.GetHelloAgent, 1)
	m[1] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](worldagent.GetWorldAgent, 1)

	floors := 3
	ts := testserver.New(m, floors)

	agents := ts.GenerateAgentArrayFromMap()

	a1 := agents[0]
	a2 := agents[1]

	messages := a1.GetAllMessages(agents)

	if len(messages) > 1 {
		t.Error("Incorrect number of messages created")
	}

	messages[0].InvokeMessageHandler(a2)
}
