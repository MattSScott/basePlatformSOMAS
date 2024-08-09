package infra

type IBaseAgent interface {
	IAgent[IBaseAgent]
}

type TestServer struct {
	IServer[IBaseAgent]
}

// func TestAgentIdOperations(t *testing.T) {
// 	var testServ infra.IServer[IBaseAgent] = TestServer{}
// 	baseAgent := infra.CreateBaseAgent[IBaseAgent](testServ)

// 	if baseAgent.GetID() == uuid.Nil {
// 		t.Error("Agent not instantiated with valid ID")
// 	}
// }

// type AgentWithState struct {
// 	*infra.BaseAgent[IBaseAgent]
// 	state int
// }

// func (aws *AgentWithState) UpdateAgentInternalState() {
// 	aws.state += 1
// }

// func TestUpdateAgentInternalState(t *testing.T) {
// 	var testServ infra.IServer[IBaseAgent] = TestServer{}

// 	ag := AgentWithState{
// 		infra.CreateBaseAgent[IBaseAgent](testServ),
// 		0,
// 	}

// 	if ag.state != 0 {
// 		t.Error("Additional agent field not correctly instantiated")
// 	}

// 	ag.UpdateAgentInternalState()

// 	if ag.state != 1 {
// 		t.Error("Agent state not correctly updated")
// 	}
// }

// func TestMessageRetrieval(t *testing.T) {
// 	var testServ infra.IServer[IBaseAgent] = TestServer{}
// 	agent := infra.CreateBaseAgent[IBaseAgent](testServ)

// 	// messages := agent.GetAllMessages([]IBaseAgent{agent})

// 	if agent == nil {
// 		t.Error("TODO!!")
// 	}

// 	// if len(messages) > 0 {
// 	// 	t.Error("Agent erroneously constructed message")
// 	// }

// }
