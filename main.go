package main

import (
	"somas_base_platform/pkg/agents/AgentTypes/agent1"
	"somas_base_platform/pkg/agents/AgentTypes/agent2"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
	infra "somas_base_platform/pkg/infra/server"
	testserver "somas_base_platform/pkg/testServer"
)

func main() {

	m := make([]infra.AgentGeneratorCountPair, 6)
	m[0] = infra.MakeAgentGeneratorCountPair(baseAgent.GetAgent, 4)
	m[1] = infra.MakeAgentGeneratorCountPair(agent2.GetAgent, 3)
	m[2] = infra.MakeAgentGeneratorCountPair(agent1.GetAgent, 2)
	floors := 3
	ts := testserver.New(m, floors)
	ts.Init()
	ts.RunGameLoop()
}
