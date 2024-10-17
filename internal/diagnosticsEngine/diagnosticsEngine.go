package diagnosticsEngine

type IDiagnosticsData interface {
	GetNumberSentMessages() int
	GetNumberMessageSuccesses() int
	GetNumberEndMessagings() int
	GetNumberMessageDrops() int
	GetMessagingSuccessRate() float32
	GetEndMessagingSuccessRate(int) float32
}

type IDiagnosticsEngine interface {
	// allow agents to report status of sent message
	ReportSendMessageStatus(bool)
	// allow server to report number of end message closures
	ReportEndMessagingStatus()
	// allow for resetting of diagnostics for round-to-round data
	ResetRoundDiagnostics()
	// compile results for end of round messaging status
	IDiagnosticsData
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

func CreateDiagnosticsEngine() *DiagnosticsEngine {
	return &DiagnosticsEngine{
		numEndMessagings:    0,
		numMessages:         0,
		numMessageSuccesses: 0,
	}
}

func (de *DiagnosticsEngine) GetNumberSentMessages() int {
	return de.numMessages
}

func (de *DiagnosticsEngine) GetNumberMessageSuccesses() int {
	return de.numMessageSuccesses
}

func (de *DiagnosticsEngine) GetNumberEndMessagings() int {
	return de.numEndMessagings
}

func (de *DiagnosticsEngine) GetNumberMessageDrops() int {
	return de.numMessages - de.numMessageSuccesses
}

func (de *DiagnosticsEngine) GetMessagingSuccessRate() float32 {
	return 100 * float32(de.numMessageSuccesses) / float32(de.numMessages)
}

func (de *DiagnosticsEngine) GetEndMessagingSuccessRate(numAgents int) float32 {
	return 100 * float32(de.numEndMessagings) / float32(numAgents)
}
