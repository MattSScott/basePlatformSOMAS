package basePlatformSOMAS

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BaseServer[T IAgent[T]] struct {
	// map of agentid -> agent struct
	agentMap map[uuid.UUID]T
	// map of agentid -> empty struct so that agents cannot access each others agent structs
	agentIdSet map[uuid.UUID]struct{}
	// map of uuid ->struct{}{} which stores the ids of agents which have stopped messaging
	agentStoppedTalkingMap map[uuid.UUID]struct{}
	//channel a server goroutine will send to in order to signal messaging completion
	messagingFinished chan struct{}
	// duration after which messaging phase forcefully ends during rounds
	turnTimeout time.Duration
	// interface which holds extended methods for round running and turn running
	roundRunner RoundRunner
	//iterations
	iterations int
	//turns
	turns int
	//lock which stops race conditions when updating agentStoppedTalkingMap
	agentMapRWMutex sync.RWMutex
	//stops multiple sends to messagingFinished during a round
	doneChannelOnce         sync.Once
	shouldRunAsyncMessaging bool
}

func (server *BaseServer[T]) HandleStartOfTurn(iter, round int) {
	server.doneChannelOnce = sync.Once{}
	server.messagingFinished = make(chan struct{})
	fmt.Printf("Iteration %d, Round %d starting...\n", iter, round)

}

func (serv *BaseServer[T]) endAgentListeningSession() bool {

	ctx, cancel := context.WithTimeout(context.Background(), serv.turnTimeout)
	defer cancel()
	select {
	case <-serv.messagingFinished:
		serv.agentStoppedTalkingMap = make(map[uuid.UUID]struct{})
		serv.shouldRunAsyncMessaging = false
		close(serv.messagingFinished)
		return true

	case <-ctx.Done():
		serv.agentStoppedTalkingMap = make(map[uuid.UUID]struct{})
		serv.shouldRunAsyncMessaging = false
		close(serv.messagingFinished)
		return false
	}
}

func (server *BaseServer[T]) HandleEndOfTurn(iter, round int) {
	server.endAgentListeningSession()
	fmt.Printf("Iteration %d, Round %d finished.\n", iter, round)
}

func (server *BaseServer[T]) SendMessage(msg IMessage[T], receivers []uuid.UUID) {
	for _, receiver := range receivers {
		go msg.InvokeMessageHandler(server.agentMap[receiver])
	}

}

func (serv *BaseServer[T]) AddAgent(agent T) {
	serv.agentMap[agent.GetID()] = agent
	serv.agentIdSet[agent.GetID()] = struct{}{}
}

func (serv *BaseServer[T]) ViewAgentIdSet() map[uuid.UUID]struct{} {
	return serv.agentIdSet
}

func (serv *BaseServer[T]) AccessAgentByID(id uuid.UUID) T {
	return serv.agentMap[id]
}

func (serv *BaseServer[T]) Initialise() {}

func (serv *BaseServer[T]) Start() {
	serv.checkHandler()
	turns := serv.turns
	iterations := serv.iterations
	for i := 0; i < iterations; i++ {
		for j := 0; j < turns; j++ {
			serv.HandleStartOfTurn(i+1, j+1)
			serv.roundRunner.RunTurn()
			serv.HandleEndOfTurn(i+1, j+1)
		}
	}
}

func (serv *BaseServer[T]) GetAgentMap() map[uuid.UUID]T {
	return serv.agentMap
}

func (serv *BaseServer[T]) agentStoppedTalking(id uuid.UUID) {
	if !serv.shouldRunAsyncMessaging {
		return
	}
	serv.agentMapRWMutex.Lock()
	serv.agentStoppedTalkingMap[id] = struct{}{}

	if len(serv.agentStoppedTalkingMap) == len(serv.agentMap) {
		serv.doneChannelOnce.Do(func() {
			serv.messagingFinished <- struct{}{}
		})
	}
	serv.agentMapRWMutex.Unlock()

}

func (serv *BaseServer[T]) SetRunHandler(handler RoundRunner) {
	serv.roundRunner = handler
}

func (serv *BaseServer[T]) checkHandler() {
	if serv.roundRunner == nil {
		panic("round running handler has not been set. Have you run SetRunHandler?")
	}
}

func (serv *BaseServer[T]) RunTurn() {}

func (bs *BaseServer[T]) RemoveAgent(agentToRemove T) {
	delete(bs.agentMap, agentToRemove.GetID())
}

func (bs *BaseServer[T]) GetIterations() int {
	return bs.iterations
}

func (bs *BaseServer[T]) RunRound() {
	//bs.roundRunner.RunRound()
}

type AgentGenerator[T IAgent[T]] func(IExposedServerFunctions[T]) T

type AgentGeneratorCountPair[T IAgent[T]] struct {
	generator AgentGenerator[T]
	count     int
}

func MakeAgentGeneratorCountPair[T IAgent[T]](generatorFunction AgentGenerator[T], count int) AgentGeneratorCountPair[T] {
	return AgentGeneratorCountPair[T]{
		generator: generatorFunction,
		count:     count,
	}
}

func (bs *BaseServer[T]) GenerateAgentArrayFromMap() []T {

	agentMapToArray := make([]T, len(bs.agentMap))

	i := 0
	for _, ag := range bs.agentMap {
		agentMapToArray[i] = ag
		i++
	}
	return agentMapToArray
}

func (bs *BaseServer[T]) SendSynchronousMessage(msg IMessage[T], recipients []uuid.UUID) {
	for _, recip := range recipients {
		if msg.GetSender() == recip {
			continue
		} else {
			msg.InvokeMessageHandler(bs.agentMap[recip])
		}
	}

}

func (bs *BaseServer[T]) RunSynchronousMessagingSession() {
	for _, agent := range bs.agentMap {
		agent.RunSynchronousMessaging()
	}
}

func (bs *BaseServer[T]) initialiseAgents(m []AgentGeneratorCountPair[T]) {

	for _, pair := range m {
		for i := 0; i < pair.count; i++ {
			agent := pair.generator(bs)
			bs.AddAgent(agent)
		}
	}

}

// generate a server instance based on a mapping function and number of iterations
func CreateServer[T IAgent[T]](generatorArray []AgentGeneratorCountPair[T], iterations, turns int, turnMaxDuration time.Duration) *BaseServer[T] {
	serv := &BaseServer[T]{
		agentMap:                make(map[uuid.UUID]T),
		agentIdSet:              make(map[uuid.UUID]struct{}),
		agentStoppedTalkingMap:  make(map[uuid.UUID]struct{}),
		turnTimeout:             turnMaxDuration,
		roundRunner:             nil,
		iterations:              iterations,
		turns:                   turns,
		messagingFinished:       make(chan struct{}),
		agentMapRWMutex:         sync.RWMutex{},
		doneChannelOnce:         sync.Once{},
		shouldRunAsyncMessaging: true,
	}
	serv.initialiseAgents(generatorArray)
	return serv
}
