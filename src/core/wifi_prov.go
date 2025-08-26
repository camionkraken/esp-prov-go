package core

import (
	"errors"
	"esp-prov-go/core/proto/protogen"
	"esp-prov-go/core/security"

	"google.golang.org/protobuf/proto"
)

func GetStatusRequest(security security.Security) ([]byte, error) {
	cfg1 := protogen.WiFiConfigPayload{
		Msg: protogen.WiFiConfigMsgType_TypeCmdGetStatus,
		Payload: &protogen.WiFiConfigPayload_CmdGetStatus{
			CmdGetStatus: &protogen.CmdGetStatus{},
		},
	}

	return EncryptMessage(security, &cfg1)
}

func GetStatusResponse(security security.Security, data []byte) (ProvStatusResult, error) {
	decMsg, err := security.DecryptData(data)
	if err != nil {
		return ProvStatusUnknown, err
	}

	var cmdResp protogen.WiFiConfigPayload
	err = proto.Unmarshal(decMsg, &cmdResp)
	if err != nil {
		return ProvStatusUnknown, err
	}

	staState := cmdResp.GetRespGetStatus().GetStaState()

	switch staState {
	case protogen.WifiStationState_Connected:
		return ProvStatusConnected, nil
	case protogen.WifiStationState_Connecting:
		return ProvStatusConnecting, nil
	case protogen.WifiStationState_Disconnected:
		return ProvStatusDisconnected, nil
	case protogen.WifiStationState_ConnectionFailed:
		failReason := cmdResp.GetRespGetStatus().GetFailReason()
		if failReason == protogen.WifiConnectFailedReason_AuthError {
			return ProvStatusAuthError, nil
		} else if failReason == protogen.WifiConnectFailedReason_NetworkNotFound {
			return ProvStatusNetworkNotFound, nil
		} else {
			return ProvStatusFailed, nil
		}
	default:
		return ProvStatusUnknown, nil
	}
}

func SetConfigRequest(security security.Security, ssid string, passphrase string) ([]byte, error) {
	cmd := protogen.WiFiConfigPayload{
		Msg: protogen.WiFiConfigMsgType_TypeCmdSetConfig,
		Payload: &protogen.WiFiConfigPayload_CmdSetConfig{
			CmdSetConfig: &protogen.CmdSetConfig{
				Ssid:       []byte(ssid),
				Passphrase: []byte(passphrase),
			},
		},
	}

	return EncryptMessage(security, &cmd)
}

func SetConfigResponse(security security.Security, data []byte) error {
	decMsg, err := security.DecryptData(data)
	if err != nil {
		return err
	}

	var resp protogen.WiFiConfigPayload
	err = proto.Unmarshal(decMsg, &resp)
	if err != nil {
		return err
	}

	if resp.GetRespSetConfig().Status != 0 {
		return errors.New("SetConfigResponse failed, status is not 0")
	}

	return nil
}

func ApplyConfigRequest(security security.Security) ([]byte, error) {
	cmd := protogen.WiFiConfigPayload{
		Msg: protogen.WiFiConfigMsgType_TypeCmdApplyConfig,
	}

	return EncryptMessage(security, &cmd)
}

func ApplyConfigResponse(security security.Security, data []byte) error {
	var resp protogen.WiFiConfigPayload
	err := proto.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	if resp.GetRespApplyConfig().Status != 0 {
		return errors.New("ApplyConfigResponse failed, status is not 0")
	}

	return nil
}
