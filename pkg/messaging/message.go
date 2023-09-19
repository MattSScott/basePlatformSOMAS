package message

// import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

type IMessage interface {
	GetSender() IAgentMessaging
	GetRecipients() []IAgentMessaging
	HowToHandleMessage(agent IAgentMessaging)
}

type BaseMessage struct {
	sender IAgentMessaging
	// content    string
	recipients []IAgentMessaging
}

// create read-only message instance
func CreateMessage(sender IAgentMessaging, recipients []IAgentMessaging) BaseMessage {
	return BaseMessage{
		sender: sender,
		// content:    content,
		recipients: recipients,
	}
}

func CreateNullMessageWithSender(sender IAgentMessaging) BaseMessage {
	return BaseMessage{
		sender: sender,
	}
}

func (bm BaseMessage) GetSender() IAgentMessaging {
	return bm.sender
}

// func (bm *BaseMessage) GetContent() string {
// 	return bm.content
// }

func (bm BaseMessage) GetRecipients() []IAgentMessaging {
	return bm.recipients
}

func (bm BaseMessage) HowToHandleMessage(agent IAgentMessaging) {
	agent.HandleMessage(bm)
}
