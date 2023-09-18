package message

//import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

type Messaging interface {
	GetMessage() Message
	HandleMessage(Message)
}
