package testserver

import (
	baseUserAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	baseServer "github.com/MattSScott/basePlatformSOMAS/pkg/infra/server"
)

type MyServer struct {
	*baseServer.BaseServer[baseUserAgent.AgentUserInterface]
	name string
}

func New(mapper []baseServer.AgentGeneratorCountPair[baseUserAgent.AgentUserInterface], numFloors int) *MyServer {
	return &MyServer{
		BaseServer: baseServer.CreateServer[baseUserAgent.AgentUserInterface](mapper, numFloors),
		name:       "TestServer",
	}
}

func (s *MyServer) GetName() string {
	return s.name
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
