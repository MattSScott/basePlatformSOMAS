package testserver

import (
	baseUserAgent "basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	baseServer "basePlatformSOMAS/pkg/infra/server"
)

type MyServer struct {
	*baseServer.BaseServer[baseUserAgent.AgentUserInterface]
	name string
}

func New(mapper []baseServer.AgentGeneratorCountPair[baseUserAgent.AgentUserInterface], numFloors int) baseServer.Server {
	return &MyServer{
		BaseServer: baseServer.CreateServer[baseUserAgent.AgentUserInterface](mapper, numFloors),
		name:       "Test Server",
	}
}

func (s *MyServer) RunGameLoop() {
	for _, agent := range s.Agents {
		agent.UpdateAgent()

	}

	for _, agent := range s.Agents {
		agent.Activity1()

	}

	for _, agent := range s.Agents {
		agent.Activity2()

	}

}
