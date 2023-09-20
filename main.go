package main

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/main/server"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
	helloagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/helloAgent"
	worldagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/worldAgent"
	testserver "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testServer"
)

func main() {

	m := make([]server.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], 5)

	m[0] = server.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](helloagent.GetHelloAgent, 3)
	m[1] = server.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](worldagent.GetWorldAgent, 2)

	iterations := 1
	ts := testserver.New(m, iterations)

	ts.Start()
}
