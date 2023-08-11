package message

import (
	

	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
)

type Letter[T baseAgent.Agent] struct {
	sender T
	message  string
	receivers    []T
}

func MessagingSession( ){

}