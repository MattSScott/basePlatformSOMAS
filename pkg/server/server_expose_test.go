package server

func (s *BaseServer[T]) ExposeStartOfTurn() {
	s.handleStartOfTurn()
}

func (s *BaseServer[T]) ExposeEndOfTurn() {
	s.handleEndOfTurn()
}

func (s *BaseServer[T]) ExposeEndListening() bool {
	return s.endAgentListeningSession()
}
