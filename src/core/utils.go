package core

import (
	"errors"
	"esp-prov-go/core/security"

	"google.golang.org/protobuf/proto"
)

func EncryptMessage(security security.Security, message proto.Message) ([]byte, error) {
	if security == nil {
		return nil, errors.New("security is nil")
	}

	if message == nil {
		return nil, errors.New("message is nil")
	}

	msgBytes, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	encMsg, err := security.EncryptData(msgBytes)
	if err != nil {
		return nil, err
	}

	return encMsg, nil
}
