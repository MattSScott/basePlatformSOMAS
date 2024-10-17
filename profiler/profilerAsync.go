package profiler

import (
	"fmt"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

type ProfilerServerAsync struct {
	*server.BaseServer[IProfilerAgent]
	startTime   time.Time
	messageGoal int32
}

func (serv *ProfilerServerAsync) RunTurn(i, j int) {
	fmt.Printf("Async Profiler Running turn %v,%v\n", i, j)
	for i := 0; i < int(serv.messageGoal); i++ {
		for _, ag := range serv.GetAgentMap() {
			newMsg := ag.CreateProfilerMessage()
			ag.BroadcastAsync(newMsg)
		}
	}
}

func (serv *ProfilerServerAsync) RunStartOfIteration(iteration int) {
	fmt.Printf("Starting iteration %v\n", iteration)
	serv.startTime = time.Now()
}

func (serv *ProfilerServerAsync) RunEndOfIteration(iteration int) {
	timeTaken := time.Since(serv.startTime)
	fmt.Printf("%v taken to run Async Messaging Round\n", timeTaken)
	fmt.Printf("Ending iteration %v\n", iteration)
}

func GenerateAsyncProfilerServer(numAgents, iterations, turns int, maxDuration time.Duration, maxThreads int, agMessagingGoal int32) *ProfilerServerAsync {
	serv := &ProfilerServerAsync{
		BaseServer:  server.CreateServer[IProfilerAgent](iterations, turns, maxDuration, maxThreads),
		messageGoal: agMessagingGoal,
	}
	for i := 0; i < numAgents; i++ {
		serv.AddAgent(NewProfilerAgent(serv, agMessagingGoal))
	}
	serv.SetGameRunner(serv)
	return serv
}
