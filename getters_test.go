package basePlatformSOMAS

import "github.com/google/uuid"

type PrivateServerFields[T IAgent[T]] interface {
	EndAgentListeningSession()
	IncrementWaitGroup()
	WaitWaitGroup()
}

func (serv *BaseServer[T]) EndAgentListeningSession() (string, bool) {
	return serv.endAgentListeningSession()
}

func (serv *BaseServer[T]) IncrementWaitGroup() {
	serv.listeningWaitGroup.Add(1)
}

func (serv *BaseServer[T]) WaitWaitGroup() {
	serv.listeningWaitGroup.Wait()
}

func (server *BaseServer[T]) MessageHandlerLimiter(msg IMessage[T], receiver uuid.UUID) {
	server.messageHandlerLimiter(msg, receiver)
}
