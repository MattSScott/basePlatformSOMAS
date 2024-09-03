package basePlatformSOMAS

type PrivateServerFields[T IAgent[T]] interface {
	EndAgentListeningSession()
	IncrementWaitGroup()
	WaitWaitGroup()
}

func (serv *BaseServer[T]) EndAgentListeningSession() bool {
	return serv.endAgentListeningSession()
}

func (serv *BaseServer[T]) EndAsyncMessaging() {
	serv.shouldRunAsyncMessaging = false
}
