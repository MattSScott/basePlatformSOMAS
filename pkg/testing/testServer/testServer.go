package testserver

import (
	baseServer "github.com/MattSScott/basePlatformSOMAS/pkg/main/server"
	"github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
)

// the base server functionality can be extended with 'environment specific' functions
type IExtendedServer interface {
	baseServer.IServer[baseExtendedAgent.IExtendedAgent]
	RunAdditionalPhase()
}

// composing the 'base server' allows access to the pre-made interface implementations
// new fields can be added to the extended server data structure
type TestServer struct {
	*baseServer.BaseServer[baseExtendedAgent.IExtendedAgent]
	name string
}

func New(mapper []baseServer.AgentGeneratorCountPair[baseExtendedAgent.IExtendedAgent], iterations int) *TestServer {
	return &TestServer{
		BaseServer: baseServer.CreateServer[baseExtendedAgent.IExtendedAgent](mapper, iterations),
		name:       "TestServer",
	}
}

func (s *TestServer) RunAdditionalPhase() {
	for _, agent := range s.GetAgents() {
		agent.GetID()
	}
}

func (s *TestServer) GetName() string {
	return s.name
}

// override to change the behaviour of the game loop
func (s *TestServer) RunGameLoop() {

	// the superclass functions can be called
	s.BaseServer.RunGameLoop()

	// able to call methods from the parametrised subclass too
	for _, agent := range s.GetAgents() {
		agent.GetPhrase()

	}

	// additional functions can be run
	s.RunAdditionalPhase()

}
