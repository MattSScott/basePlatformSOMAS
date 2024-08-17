package infra

import (
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
	// waitgroup for message synchronisation
	listeningWaitGroup *sync.WaitGroup
	// map of uuid ->struct{}{} which stores the ids of agents which have stopped messaging
	agentStoppedTalkingMap map[uuid.UUID]struct{}
	// duration after which messaging phase forcefully ends during rounds
	maxMessagingDuration time.Duration
	// run single round
	roundRunner RoundRunner
	//iterations
	iterations int
}

func (server *BaseServer[T]) HandleStartOfTurn(iter, round int) {
	fmt.Printf("Iteration %d, Round %d starting...\n", iter, round)

}

func (serv *BaseServer[T]) endAgentListeningSession() {

	timeoutChannel := time.After(serv.maxMessagingDuration)

agentMessaging:
	for {
		select {
		case <-timeoutChannel:
			fmt.Println("len of stoppedtalkingmap:,", len(serv.agentStoppedTalkingMap))
			fmt.Println("Stopped messaging at time limit", serv.maxMessagingDuration, "seconds")
			break agentMessaging

		default:
			if len(serv.agentStoppedTalkingMap) == len(serv.agentMap) {
				fmt.Println("stopped messaging early", len(serv.agentStoppedTalkingMap), len(serv.agentMap))
				break agentMessaging
			}

		}
	}
	serv.agentStoppedTalkingMap = make(map[uuid.UUID]struct{})
}

func (server *BaseServer[T]) HandleEndOfTurn(iter, round int) {
	server.endAgentListeningSession()
	fmt.Printf("Iteration %d, Round %d finished.\n", iter, round)
}

func (server *BaseServer[T]) RunAgentLoop() {}

func (server *BaseServer[T]) SendMessage(msg IMessage[T], receivers []uuid.UUID) {
	//defer server.listeningWaitGroup.Done()
	for _, receiver := range receivers {
		msg.InvokeMessageHandler(server.agentMap[receiver])
	}

}

func (serv *BaseServer[T]) ViewAgentMap() map[uuid.UUID]T {
	return serv.agentMap
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
	turns := 5
	iterations := 1
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
	fmt.Println("sending stop talking request,id:", id)
	// select {
	// case serv.agentServerChannel <- id:
	// default:
	// }
	serv.agentStoppedTalkingMap[id] = struct{}{}
}

func (serv *BaseServer[T]) SetRunHandler(handler RoundRunner) {
	serv.roundRunner = handler
}

func (serv *BaseServer[T]) checkHandler() {
	if serv.roundRunner == nil {
		panic("handler has not been set. Have you run SetRunHandler?")
	}
}

func (serv *BaseServer[T]) RunTurn() {}

// func (bs *BaseServer[T]) GetAgentMap() map[uuid.UUID]T {
// 	return bs.agentMap
// }

// func (bs *BaseServer[T]) AddAgent(agentToAdd T) {
// 	bs.agentMap[agentToAdd.GetID()] = agentToAdd
// }

func (bs *BaseServer[T]) RemoveAgent(agentToAdd T) {
	delete(bs.agentMap, agentToAdd.GetID())
}

func (bs *BaseServer[T]) GetIterations() int {
	return bs.iterations
}

func (bs *BaseServer[T]) RunGameLoop() {
	if bs.roundRunner == nil {
		panic("roundRunner has not been set.")

	}
	for id, agent := range bs.agentMap {
		fmt.Printf("Agent %s updating state \n", id)
		agent.UpdateAgentInternalState()
	}
}

func (bs *BaseServer[T]) RunRound() {}

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
		}
		msg.InvokeMessageHandler(bs.agentMap[recip])
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
func CreateServer[T IAgent[T]](generatorArray []AgentGeneratorCountPair[T], iterations int, maxDuration time.Duration) *BaseServer[T] {
	serv := &BaseServer[T]{
		agentMap:               make(map[uuid.UUID]T),
		agentIdSet:             make(map[uuid.UUID]struct{}),
		listeningWaitGroup:     &sync.WaitGroup{},
		agentStoppedTalkingMap: make(map[uuid.UUID]struct{}),
		maxMessagingDuration:   maxDuration,
		roundRunner:            nil,
		iterations:             iterations,
	}
	serv.initialiseAgents(generatorArray)
	return serv
}
