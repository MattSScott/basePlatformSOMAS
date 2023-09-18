package message

//import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

type Messaging[T Agent] interface {
	GetMessage() Message[T]
	HandleMessage()
}
