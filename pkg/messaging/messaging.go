package message

import "github.com/google/uuid"

//import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

type IAgentMessaging interface {
	GetID() uuid.UUID
	GetMessage() Message
	HandleMessage(Message)
}
