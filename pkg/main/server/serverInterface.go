package server

import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseAgent"

type IServer[T baseagent.IAgent[T]] interface {
	// the set of functions defining how a 'game loop' should run
	RunGameLoop()
	// begins simulator
	Start()
	// gives access to the agents in the simulator
	GetAgents() []T
	// gives access to number of iteration in simulator
	GetNumTurns() int
}
