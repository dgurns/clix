package session

type State int

const (
	StateSystemAsking State = iota
	StateUserAnswering
	StateToolUsing
)

type Session struct {
	State State
}
