package security

import (
	"errors"
	"esp-prov-go/core/proto/protogen"

	"google.golang.org/protobuf/proto"
)

type Security0 struct {
	sessionState int
}

func (s *Security0) SecuritySession(data []byte) ([]byte, error) {
	switch s.sessionState {
	case 0:
		s.sessionState = 1
		req, err := s.setup0Request()

		if err != nil {
			return nil, err
		}

		return req, nil
	case 1:
		if data == nil {
			return nil, errors.New("data should not be nil")
		}

		err := s.setup0Response(data)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *Security0) setup0Request() ([]byte, error) {
	setupReq := protogen.SessionData{
		SecVer: 0,
		Proto: &protogen.SessionData_Sec0{
			Sec0: &protogen.Sec0Payload{
				Payload: &protogen.Sec0Payload_Sc{
					Sc: &protogen.S0SessionCmd{},
				},
			},
		},
	}

	bytes, err := proto.Marshal(&setupReq)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (s *Security0) setup0Response(responseData []byte) error {
	setupResp := protogen.SessionData{}
	err := proto.Unmarshal(responseData, &setupResp)
	return err
}

func (s *Security0) DecryptData(data []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Security0) EncryptData(bytes []byte) {
	//TODO implement me
	panic("implement me")
}
