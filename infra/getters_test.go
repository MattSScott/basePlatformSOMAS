package infra

type PrivateServerFields[T IAgent[T]] interface {
	EndAgentListeningSession()
	IncrementWaitGroup()
	WaitWaitGroup()
}

func (serv *BaseServer[T]) EndAgentListeningSession() {
	serv.endAgentListeningSession()
}

func (serv *BaseServer[T]) IncrementWaitGroup() {
	serv.listeningWaitGroup.Add(1)
}

func (serv *BaseServer[T]) WaitWaitGroup() {
	serv.listeningWaitGroup.Wait()
}
