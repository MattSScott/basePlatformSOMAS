package messaging

import (
	baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
)

func distributeMessages(message message.Message, recipients []message.IAgentMessaging) {
	for _, recip := range recipients {
		recip.HandleMessage(message)
	}
}

func MessagingSession(agents []baseAgent.IAgent) {

	for _, agent := range agents {
		messageFromAgent := agent.GetMessage()
		distributeMessages(messageFromAgent, messageFromAgent.GetRecipients())
	}

}
