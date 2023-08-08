package serverTesting

import (
	//agent1 "basePlatformSOMAS/pkg/agents/AgentTesting/agent1"
	//agent2 "basePlatformSOMAS/pkg/agents/AgentTesting/agent2"
	//baseUserAgent "basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
	infra "basePlatformSOMAS/pkg/infra/server"
	//testserver "basePlatformSOMAS/pkg/testServer"
	"github.com/google/uuid"
	"testing"
)


func TestMakeServerBase(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[baseAgent.Agent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(baseAgent.GetAgent, 4)

	serv := infra.CreateServer[baseAgent.Agent](m, 5)
	

	for _, agent := range serv.Agents {

		if agent.GetID() == uuid.Nil {
			t.Error("Error creating agent")

		}
		
	}
	serv.RunGameLoop()
	serv.Start()

}

// func TestmakeServerTest() {

// 	m := make([]infra.AgentGeneratorCountPair[baseUserAgent.AgentUserInterface], 2)
// 	m[0] = infra.MakeAgentGeneratorCountPair[baseUserAgent.AgentUserInterface](agent2.GetAgent, 3)
// 	m[1] = infra.MakeAgentGeneratorCountPair[baseUserAgent.AgentUserInterface](agent1.GetAgent, 2)
// 	floors := 3
// 	ts := testserver.New(m, floors)
// 	ts.RunGameLoop()
// 	ts.Start()
// }