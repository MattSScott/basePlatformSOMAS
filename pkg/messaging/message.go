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

func CreateLetter(){

}

type LetterI [T baseAgent.Agent]interface { //better name would be great
	GetSnd() T
	SetSnd(a T)
	GetMsg() string
	SetMsg(s string)

} 

func (l *Letter[T])GetSnd() T {
	return l.sender
}
func (l *Letter[T])SetSnd(a T) {
	l.sender= a
}

func (l *Letter[T])GetMsg() string{
	return l.message
}

func (l *Letter[T])SetMsg(s string){
	l.message = s
} 

func (l *Letter[T])GetRcv() []T{
	return l.receivers
}

func (l *Letter[T])SetRcv(a []T) {
	l.receivers=a
}