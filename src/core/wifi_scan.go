package core

import (
	"errors"
	"esp-prov-go/core/proto/protogen"
	"esp-prov-go/core/security"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func StartScanRequest(
	security security.Security,
	blocking bool, passive bool,
	groupChannels uint32,
	periodMs uint32) ([]byte, error) {

	cmd := protogen.WiFiScanPayload{
		Msg: protogen.WiFiScanMsgType_TypeCmdScanStart,
		Payload: &protogen.WiFiScanPayload_CmdScanStart{
			CmdScanStart: &protogen.CmdScanStart{
				Blocking:      blocking,
				Passive:       passive,
				GroupChannels: groupChannels,
				PeriodMs:      periodMs,
			},
		},
	}

	return EncryptMessage(security, &cmd)
}

func StartScanResponse(security security.Security, data []byte) error {
	decResp, err := security.DecryptData(data)
	if err != nil {
		return err
	}

	var resp protogen.WiFiScanPayload
	err = proto.Unmarshal(decResp, &resp)

	if err != nil {
		return err
	}

	if resp.Status != 0 {
		return errors.New(fmt.Sprintf("Error in Start Scan Response: Status %d", resp.Status))
	}

	return nil
}

func ScanStatusRequest(security security.Security) ([]byte, error) {
	cmd := protogen.WiFiScanPayload{
		Msg: protogen.WiFiScanMsgType_TypeCmdScanStatus,
	}

	return EncryptMessage(security, &cmd)
}

func ScanStatusResponse(security security.Security, data []byte) (ScanStatusResult, error) {
	decResp, err := security.DecryptData(data)
	if err != nil {
		return ScanStatusResult{}, err
	}

	var resp protogen.WiFiScanPayload
	err = proto.Unmarshal(decResp, &resp)
	if err != nil {
		return ScanStatusResult{}, err
	}

	if resp.Status != 0 {
		return ScanStatusResult{}, errors.New(fmt.Sprintf("Error in Scan Status Response: Status %d", resp.Status))
	}

	return ScanStatusResult{
		Count:    resp.GetRespScanStatus().ResultCount,
		Finished: resp.GetRespScanStatus().ScanFinished,
	}, nil
}

func ScanResultRequest(security security.Security, index uint32, count uint32) ([]byte, error) {
	cmd := protogen.WiFiScanPayload{
		Msg: protogen.WiFiScanMsgType_TypeCmdScanResult,
		Payload: &protogen.WiFiScanPayload_CmdScanResult{
			CmdScanResult: &protogen.CmdScanResult{
				StartIndex: index,
				Count:      count,
			},
		},
	}

	return EncryptMessage(security, &cmd)
}

func ScanResultResponse(security security.Security, data []byte) ([]WiFiFromScan, error) {
	decResp, err := security.DecryptData(data)
	if err != nil {
		return nil, err
	}

	var resp protogen.WiFiScanPayload
	err = proto.Unmarshal(decResp, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Status != 0 {
		return nil, errors.New(fmt.Sprintf("Error in Scan Result Response: Status %d", resp.Status))
	}

	wiFis := make([]WiFiFromScan, 0)

	for _, wiFi := range resp.GetRespScanResult().Entries {
		wiFiScanResult, err := FromWiFiScanResult(wiFi)
		if err != nil {
			continue
		}
		wiFis = append(wiFis, wiFiScanResult)
	}

	return wiFis, nil
}
