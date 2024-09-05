package basePlatformSOMAS

import (
	"sync"
	"sync/atomic"
)

type PrivateServerFields[T IAgent[T]] interface {
	EndAgentListeningSession()
}

func (serv *BaseServer[T]) EndAgentListeningSession() bool {
	return serv.endAgentListeningSession()
}

func (serv *BaseServer[T]) EndAsyncMessaging() {
	serv.shouldAllowStopTalking = false
}

func (ag *BaseAgent[T]) NotifyAgentFinishedMessagingUnthreaded(wg *sync.WaitGroup, counter *uint32) {
	defer wg.Done()
	ag.agentStoppedTalking(ag.id)
	atomic.AddUint32(counter, 1)
}
