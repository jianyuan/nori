package message

//go:generate stringer -type=State

type State int

const (
	Success State = iota
	Failure
	Revoked
	Started
	Received
	Retry
	Pending
)
