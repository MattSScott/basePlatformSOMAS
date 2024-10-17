package main

import (
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/profiler"
)

func main() {
	//benchmarks
	agents := 2
	rounds := 1
	iterations := 1
	timeLimit := time.Second
	threads := 100
	
	syncProfiler := profiler.GenerateSyncProfilerServer(agents, rounds, iterations, timeLimit, threads,2)
	asyncProfiler := profiler.GenerateAsyncProfilerServer(agents, rounds, iterations, timeLimit, threads,2)
	syncProfiler.Start()
	asyncProfiler.Start()

}
