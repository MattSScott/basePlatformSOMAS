// package main

// import (
// 	testserver "somas_base_platform/pkg/testServer"
// )

// func main() {
// 	ts := testserver.New()
// 	ts.Init()
// 	ts.Start()
// }

package main

import (
	agent1 "somas_base_platform/pkg/agents/AgentTypes/agent1"
	agent2 "somas_base_platform/pkg/agents/AgentTypes/agent2"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
)

func main() {

	var Agents []baseAgent.Agent = make([]baseAgent.Agent, 3)
	Agents[0] = agent1.GetAgent()
	Agents[1] = agent1.GetAgent()
	Agents[2] = agent2.GetAgent()

	for _, ag := range Agents {
		ag.Activity()
	}

}
