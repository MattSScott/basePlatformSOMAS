package message

import baseagent "basePlatformSOMAS/pkg/agents/BaseAgent"

type Messaging[T baseagent.Agent] interface {
	GetMessage() Message[T]
	HandleMessage()
}
