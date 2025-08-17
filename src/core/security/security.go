package security

type Security interface {
	SecuritySession(data []byte) ([]byte, error)
	EncryptData(data []byte) ([]byte, error)
	DecryptData(data []byte) ([]byte, error)
}
