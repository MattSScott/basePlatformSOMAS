package diagnosticsEngine

import "fmt"

type IDiagnosticsEngine interface {
	// allow agents to report status of sent message
	ReportSendMessageStatus(bool)
	// allow server to report number of end message closures
	ReportEndMessagingStatus()
	// allow for resetting of diagnostics for round-to-round data
	ResetRoundDiagnostics()
	// compile results for end of round messaging status
	CompileRoundDiagnostics(int)
}

type DiagnosticsEngine struct {
	numEndMessagings    int
	numMessages         int
	numMessageSuccesses int
}

func (de *DiagnosticsEngine) ReportSendMessageStatus(status bool) {
	de.numMessages++
	if status {
		de.numMessageSuccesses++
	}
}

func (de *DiagnosticsEngine) ReportEndMessagingStatus() {
	de.numEndMessagings++
}

func (de *DiagnosticsEngine) ResetRoundDiagnostics() {
	de.numEndMessagings = 0
	de.numMessages = 0
	de.numMessageSuccesses = 0
}

func (de *DiagnosticsEngine) CompileRoundDiagnostics(numAgents int) {
	messageSuccessRate := float32(de.numMessages) / float32(de.numMessageSuccesses)
	msgDropped := de.numMessages - de.numMessageSuccesses
	endMessagingSuccessRate := float32(de.numEndMessagings) / float32(numAgents)
	fmt.Printf("%f%% of messages successfully sent (%d delivered, %d dropped)\n", messageSuccessRate, de.numMessageSuccesses, msgDropped)
	fmt.Printf("%f%% of agents successfully ended messaging (%d ended, %d total)\n", endMessagingSuccessRate, de.numEndMessagings, numAgents)
}

func CreateDiagnosticsEngine() *DiagnosticsEngine {
	return &DiagnosticsEngine{
		numEndMessagings:    0,
		numMessages:         0,
		numMessageSuccesses: 0,
	}
}
