package testserver

import (
	baseServer "github.com/MattSScott/basePlatformSOMAS/pkg/main/server"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
)

type MyServer struct {
	*baseServer.BaseServer[baseExtendedAgent.IExtendedAgent]
	name string
}

func New(mapper []baseServer.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], iterations int) *MyServer {
	return &MyServer{
		BaseServer: baseServer.CreateServer[baseExtendedAgent.IExtendedAgent](mapper, iterations),
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

	// able to call methods from the parametrised subclass
	for _, agent := range s.Agents {
		agent.GetPhrase()

	}

	s.RunMessagingSession()

}
