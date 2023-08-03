package testserver

import (
	"fmt"
	agent1 "somas_base_platform/pkg/agents/AgentTypes/agent1"
	agent2 "somas_base_platform/pkg/agents/AgentTypes/agent2"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
	baseServer "somas_base_platform/pkg/infra/server"
	infra "somas_base_platform/pkg/infra/server"
)

type MyServer struct {
	// // embed struct for composition
	// baseServer.BaseServer
	// keep second reference to superclass for unimplemented methods
	*baseServer.BaseServer
	name      string
	numFloors int
}

func New() MyServer {
	return MyServer{
		BaseServer: &baseServer.BaseServer{},
		numFloors:  20,
	}
}

func (ms *MyServer) Init() {
	ms.name = "My Server"
	fmt.Printf("Name field added as: %s \n", ms.name)
	m := make([]infra.AgentGeneratorCountPair, 6)
	m[0] = infra.MakeAgentGeneratorCountPair(baseAgent.GetAgent, 4)
	m[1] = infra.MakeAgentGeneratorCountPair(agent2.GetAgent, 3)
	m[2] = infra.MakeAgentGeneratorCountPair(agent1.GetAgent, 2)

	ms.BaseServer.Init(m)
}
