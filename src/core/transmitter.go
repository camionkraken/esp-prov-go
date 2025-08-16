package core

type Transmitter interface {
	Send(path string, data []byte) ([]byte, error)
}
