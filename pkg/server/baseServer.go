package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/google/uuid"
)

type BaseServer[T agent.IAgent[T]] struct {
	// map of agentid -> agent struct
	agentMap map[uuid.UUID]T
	// map of agentid -> empty struct so that agents cannot access each others agent structs
	agentIdSet map[uuid.UUID]struct{}
	// channel a server goroutine will send to in order to signal messaging completion
	agentFinishedMessaging chan uuid.UUID
	// duration after which messaging phase forcefully ends during rounds
	turnTimeout time.Duration
	// interface which holds extended methods for round running and turn running
	roundRunner RoundRunner
	// number of iterations for server
	iterations int
	// number of turns for server
	turns int
	// mutex for agentStoppedTalkingMap access
	endNotifyAgentDone endNotifyAgentDone
	// stops multiple sends to messagingFinished during a round
	doneChannelOnce sync.Once
	// flag to disable async message propagation after timeout
	shouldAllowStopTalking bool
}

type endNotifyAgentDone struct {
	endNotifyAgentDoneContext context.Context
	cancelNotifyAgentDone     context.CancelFunc
}

func (server *BaseServer[T]) HandleStartOfTurn(iter, round int) {
	server.doneChannelOnce = sync.Once{}
	server.agentFinishedMessaging = make(chan uuid.UUID)

	fmt.Printf("Iteration %d, Round %d starting...\n", iter, round)
}

func (serv *BaseServer[T]) resetServerAsyncHelpers() {

	serv.endNotifyAgentDone.cancelNotifyAgentDone()
	serv.shouldAllowStopTalking = false
	newCtx, newCancel := context.WithCancel(context.Background())
	serv.endNotifyAgentDone.endNotifyAgentDoneContext = newCtx
	serv.endNotifyAgentDone.cancelNotifyAgentDone = newCancel
	close(serv.agentFinishedMessaging)

}

func (serv *BaseServer[T]) EndAgentListeningSession() bool {
	status := true
	ctx, cancel := context.WithTimeout(context.Background(), serv.turnTimeout)
	defer cancel()

	agentStoppedTalkingMap := make(map[uuid.UUID]struct{})
awaitSessionEnd:
	for len(agentStoppedTalkingMap) != len(serv.agentMap) {
		select {
		case id := <-serv.agentFinishedMessaging:
			agentStoppedTalkingMap[id] = struct{}{}

		case <-ctx.Done():
			status = false
			break awaitSessionEnd
		}
	}
	serv.resetServerAsyncHelpers()
	return status
}

func (server *BaseServer[T]) HandleEndOfTurn(iter, round int) {
	server.EndAgentListeningSession()
	fmt.Printf("Iteration %d, Round %d finished.\n", iter, round)
}

func (server *BaseServer[T]) SendMessage(msg message.IMessage[T], receivers []uuid.UUID) {
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
	for i := 0; i < serv.iterations; i++ {
		for j := 0; j < serv.turns; j++ {
			serv.HandleStartOfTurn(i+1, j+1)
			serv.roundRunner.RunTurn()
			serv.HandleEndOfTurn(i+1, j+1)
		}
	}
}

func (serv *BaseServer[T]) GetAgentMap() map[uuid.UUID]T {
	return serv.agentMap
}

func (serv *BaseServer[T]) AgentStoppedTalking(id uuid.UUID) {
	if !serv.shouldAllowStopTalking {
		return
	}
	select {
	case serv.agentFinishedMessaging <- id:
	case <-serv.endNotifyAgentDone.endNotifyAgentDoneContext.Done():
	}
}

func (serv *BaseServer[T]) SetRunHandler(handler RoundRunner) {
	serv.roundRunner = handler
}

func (serv *BaseServer[T]) checkHandler() {
	if serv.roundRunner == nil {
		panic("round running handler has not been set. Have you run SetRunHandler?")
	}
}

func (serv *BaseServer[T]) RunTurn() {
	serv.roundRunner.RunTurn()
}

func (serv *BaseServer[T]) RemoveAgent(agentToRemove T) {
	delete(serv.agentMap, agentToRemove.GetID())
	delete(serv.agentIdSet, agentToRemove.GetID())
}

func (serv *BaseServer[T]) GetIterations() int {
	return serv.iterations
}

func (serv *BaseServer[T]) RunRound() {
	serv.roundRunner.RunRound()
}

func (serv *BaseServer[T]) GenerateAgentArrayFromMap() []T {
	agentMapToArray := make([]T, len(serv.agentMap))

	i := 0
	for _, ag := range serv.agentMap {
		agentMapToArray[i] = ag
		i++
	}
	return agentMapToArray
}

func (serv *BaseServer[T]) SendSynchronousMessage(msg message.IMessage[T], recipients []uuid.UUID) {
	for _, recip := range recipients {
		if msg.GetSender() == recip {
			continue
		}
		msg.InvokeMessageHandler(serv.agentMap[recip])
	}
}

func (serv *BaseServer[T]) RunSynchronousMessagingSession() {
	serv.shouldAllowStopTalking = false
	for _, agent := range serv.agentMap {
		agent.RunSynchronousMessaging()
	}
	serv.shouldAllowStopTalking = true
}

func (serv *BaseServer[T]) initialiseAgents(m []agent.AgentGeneratorCountPair[T]) {
	for _, pair := range m {
		for i := 0; i < pair.Count; i++ {
			agent := pair.Generator(serv)
			serv.AddAgent(agent)
		}
	}
}

// generate a server instance based on a mapping function and number of iterations
func CreateServer[T agent.IAgent[T]](generatorArray []agent.AgentGeneratorCountPair[T], iterations, turns int, turnMaxDuration time.Duration) *BaseServer[T] {
	serv := &BaseServer[T]{
		agentMap:               make(map[uuid.UUID]T),
		agentIdSet:             make(map[uuid.UUID]struct{}),
		turnTimeout:            turnMaxDuration,
		roundRunner:            nil,
		iterations:             iterations,
		turns:                  turns,
		agentFinishedMessaging: make(chan uuid.UUID),
		endNotifyAgentDone:     endNotifyAgentDone{},
		doneChannelOnce:        sync.Once{},
		shouldAllowStopTalking: true,
	}
	ctx, cancel := context.WithCancel(context.Background())
	serv.endNotifyAgentDone.endNotifyAgentDoneContext = ctx
	serv.endNotifyAgentDone.cancelNotifyAgentDone = cancel
	serv.initialiseAgents(generatorArray)
	return serv
}

func (serv *BaseServer[T]) EndAsyncMessaging() {
	serv.shouldAllowStopTalking = false
}
