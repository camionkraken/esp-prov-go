package core

type ProvStatusResult int

const (
	ProvStatusUnknown ProvStatusResult = iota
	ProvStatusConnected
	ProvStatusFailed
	ProvStatusNetworkNotFound
	ProvStatusAuthError
	ProvStatusDisconnected
	ProvStatusConnecting
)
