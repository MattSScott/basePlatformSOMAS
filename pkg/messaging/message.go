package message

// import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

// base interface structure used for message passing - can be composed for more complex message structures

// new message types extend this
type IMessage interface {
	GetSender() IAgentMessenger
	GetRecipients() []IAgentMessenger
	Accept(IAgentMessenger)
}

type BaseMessage struct {
	sender     IAgentMessenger
	recipients []IAgentMessenger
}

// create read-only message instance
func CreateMessage(sender IAgentMessenger, recipients []IAgentMessenger) BaseMessage {
	return BaseMessage{
		sender:     sender,
		recipients: recipients,
	}
}

func CreateNullMessageWithSender(sender IAgentMessenger) BaseMessage {
	return BaseMessage{
		sender: sender,
	}
}

func (bm BaseMessage) GetSender() IAgentMessenger {
	return bm.sender
}

func (bm BaseMessage) GetRecipients() []IAgentMessenger {
	return bm.recipients
}

func (bm BaseMessage) Accept(agent IAgentMessenger) {
}
