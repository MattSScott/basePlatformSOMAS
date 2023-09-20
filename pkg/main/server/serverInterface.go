package server

type IServer interface {
	// the set of functions defining a 'game loop' should run
	RunGameLoop()
	// runs simulator
	Start()
}
