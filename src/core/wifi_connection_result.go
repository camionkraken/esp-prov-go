package core

type WiFiConnectionResult int

const (
	WiFiConnected WiFiConnectionResult = iota
	WiFiNetworkNotFound
	WiFiAuthError
	WiFiRetriesExceeded
	WiFiSetConfigurationFailed
	WiFiApplyConfigurationFailed
	WiFiFailed
)
