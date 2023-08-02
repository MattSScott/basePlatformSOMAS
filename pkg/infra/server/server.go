package infra

type Server interface {
	// initial environment variables should be instantiated
	Init()
	// the set of functions defining a 'game loop' should run
	RunGameLoop()
	// runs simulator
	Start()
}
