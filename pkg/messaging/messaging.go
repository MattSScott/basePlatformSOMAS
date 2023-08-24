package message

import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

type Messaging[T baseagent.Agent] interface {
	GetMessage() Message[T]
	HandleMessage()
}
