package main

import (
	baseserver "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseServer"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
	helloagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/helloAgent"
	worldagent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/worldAgent"
	testserver "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testServer"
)

func main() {

	m := make([]baseserver.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], 2)

	m[0] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](helloagent.GetHelloAgent, 3)
	m[1] = baseserver.MakeAgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent](worldagent.GetWorldAgent, 2)

	iterations := 1
	ts := testserver.New(m, iterations)

	ts.Start()
}
