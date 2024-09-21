package server

import (
	"context"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/google/uuid"
)

type BaseServer[T agent.IAgent[T]] struct {
	// map of agentid -> agent struct
	agentMap map[uuid.UUID]T
	// hashset of agent IDs
	agentIdSet map[uuid.UUID]struct{}
	// channel a server goroutine will send to in order to signal messaging completion
	agentFinishedMessaging chan uuid.UUID
	// duration after which messaging phase forcefully ends during turns
	turnTimeout time.Duration
	// interface which allows overridable turns
	gameRunner GameRunner
	// number of iterations for server
	iterations int
	// number of turns for server
	turns int
	// closable channel to signify that messaging is complete
	endNotifyAgentDone chan struct{}
	// limits the number of sendmessage goroutines executing at once
	messageSenderSemaphore chan struct{}
}

func (server *BaseServer[T]) HandleStartOfTurn() {
	server.agentFinishedMessaging = make(chan uuid.UUID)
	server.endNotifyAgentDone = make(chan struct{})
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
	close(serv.endNotifyAgentDone)
	return status
}

func (server *BaseServer[T]) HandleEndOfTurn() {
	server.EndAgentListeningSession()
}

func (server *BaseServer[T]) SendMessage(msg message.IMessage[T], recipient uuid.UUID) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	select {
	case server.messageSenderSemaphore <- struct{}{}:
		go func() {
			msg.InvokeMessageHandler(server.agentMap[recipient])
			<-server.messageSenderSemaphore
		}()
	default:
	}
}

func (server *BaseServer[T]) BroadcastMessage(msg message.IMessage[T]) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	for id := range server.ViewAgentIdSet() {
		if id == msg.GetSender() {
			continue
		}
		server.SendMessage(msg, id)
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

func (serv *BaseServer[T]) Start() {
	serv.checkHandler()
	for i := 0; i < serv.iterations; i++ {
		serv.RunStartOfIteration(i)
		for j := 0; j < serv.turns; j++ {
			serv.HandleStartOfTurn()
			serv.RunTurn(i, j)
			serv.HandleEndOfTurn()
		}
		serv.RunEndOfIteration(i)
	}
}

func (serv *BaseServer[T]) GetAgentMap() map[uuid.UUID]T {
	return serv.agentMap
}

func (serv *BaseServer[T]) AgentStoppedTalking(id uuid.UUID) {
	select {
	case serv.agentFinishedMessaging <- id:
		return
	case <-serv.endNotifyAgentDone:
		return
	}
}

func (serv *BaseServer[T]) SetTurnRunner(handler GameRunner) {
	serv.gameRunner = handler
}

func (serv *BaseServer[T]) checkHandler() {
	if serv.gameRunner == nil {
		panic("Handler for running turn has not been set. Have you called SetTurnRunner?")
	}
}

func (serv *BaseServer[T]) RunTurn(turn, iteration int) {
	serv.gameRunner.RunTurn(turn, iteration)
}

func (serv *BaseServer[T]) RunStartOfIteration(iteration int) {
	serv.gameRunner.RunStartOfIteration(iteration)
}

func (serv *BaseServer[T]) RunEndOfIteration(iteration int) {
	serv.gameRunner.RunEndOfIteration(iteration)
}

func (serv *BaseServer[T]) GetTurns() int {
	return serv.turns
}

func (serv *BaseServer[T]) GetIterations() int {
	return serv.iterations
}

func (serv *BaseServer[T]) RemoveAgent(agentToRemove T) {
	delete(serv.agentMap, agentToRemove.GetID())
	delete(serv.agentIdSet, agentToRemove.GetID())
}

func (serv *BaseServer[T]) SendSynchronousMessage(msg message.IMessage[T], recip uuid.UUID) {
	if msg.GetSender() == uuid.Nil {
		panic("No sender found - did you compose the BaseMessage?")
	}
	if msg.GetSender() == recip {
		return
	}
	msg.InvokeMessageHandler(serv.agentMap[recip])
}

func (serv *BaseServer[T]) RunSynchronousMessagingSession() {
	for _, agent := range serv.agentMap {
		agent.RunSynchronousMessaging()
	}
}

// generate a server instance based on a mapping function and number of iterations
func CreateServer[T agent.IAgent[T]](iterations, turns int, turnMaxDuration time.Duration, messageBandwidth int) *BaseServer[T] {
	return &BaseServer[T]{
		agentMap:               make(map[uuid.UUID]T),
		agentIdSet:             make(map[uuid.UUID]struct{}),
		turnTimeout:            turnMaxDuration,
		gameRunner:             nil,
		iterations:             iterations,
		turns:                  turns,
		agentFinishedMessaging: make(chan uuid.UUID),
		endNotifyAgentDone:     make(chan struct{}),
		messageSenderSemaphore: make(chan struct{}, messageBandwidth),
	}
}
