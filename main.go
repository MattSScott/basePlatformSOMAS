package main

import (
	"fmt"

	agent1 "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/agent1"
	agent2 "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/agent2"
	baseUserAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
	infra "github.com/MattSScott/basePlatformSOMAS/pkg/infra/server"
	testserver "github.com/MattSScott/basePlatformSOMAS/pkg/testServer"
)

func main() {

	fmt.Println("Running base server:")
	makeServerBase()
	fmt.Println("Base server finished.")

	fmt.Println("----------------------------")

	fmt.Println("Running test server:")
	makeServerTest()
	fmt.Println("Test server finished.")

}

func makeServerBase() {
	m := make([]infra.AgentGeneratorCountPair[baseAgent.IAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(baseAgent.GetAgent, 4)

	serv := infra.CreateServer[baseAgent.IAgent](m, 5)

	serv.RunGameLoop()
	serv.Start()

}

func makeServerTest() {

	m := make([]infra.AgentGeneratorCountPair[baseUserAgent.AgentUserInterface], 2)
	m[0] = infra.MakeAgentGeneratorCountPair[baseUserAgent.AgentUserInterface](agent2.GetAgent, 3)
	m[1] = infra.MakeAgentGeneratorCountPair[baseUserAgent.AgentUserInterface](agent1.GetAgent, 2)
	floors := 3
	ts := testserver.New(m, floors)
	ts.RunGameLoop()
	ts.Start()
}
