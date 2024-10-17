package profiler

import (
	"fmt"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

type ProfilerServerSync struct {
	*server.BaseServer[IProfilerAgent]
	startTime   time.Time
	messageGoal int32
}

func (serv *ProfilerServerSync) RunTurn(i, j int) {
	fmt.Printf("Sync Profiler Running turn %v,%v\n", i, j)
	for i := 0; i < int(serv.messageGoal); i++ {
		for _, ag := range serv.GetAgentMap() {
			newMsg := ag.CreateProfilerMessage()
			ag.BroadcastSync(newMsg)
		}
	}
}

func (serv *ProfilerServerSync) RunStartOfIteration(iteration int) {
	fmt.Printf("Starting iteration %v\n", iteration)
	serv.startTime = time.Now()
}

func (serv *ProfilerServerSync) RunEndOfIteration(iteration int) {
	timeTaken := time.Since(serv.startTime)
	fmt.Printf("%v taken to run Sync Messaging Round\n", timeTaken)
	fmt.Printf("Ending iteration %v\n", iteration)
}

func GenerateSyncProfilerServer(numAgents, iterations, turns int, maxDuration time.Duration, maxThreads int, agMessageGoal int32) *ProfilerServerSync {
	serv := &ProfilerServerSync{
		BaseServer:  server.CreateServer[IProfilerAgent](iterations, turns, maxDuration, maxThreads),
		messageGoal: agMessageGoal,
	}
	for i := 0; i < numAgents; i++ {
		serv.AddAgent(NewProfilerAgent(serv, agMessageGoal))
	}
	serv.SetGameRunner(serv)
	return serv
}
