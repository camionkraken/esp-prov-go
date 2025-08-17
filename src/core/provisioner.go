package core

import (
	"encoding/json"
	"esp-prov-go/core/security"
)

type Provisioner struct {
	transmitter Transmitter
	security    security.Security
}

func (provisioner *Provisioner) getTransmitter() Transmitter {
	return provisioner.transmitter
}

func (provisioner *Provisioner) getSecurity() security.Security {
	return provisioner.security
}

func NewProvisioner(transmitter Transmitter, security security.Security) *Provisioner {
	return &Provisioner{
		transmitter: transmitter,
		security:    security,
	}
}

func (provisioner *Provisioner) GetProtoVersion() (ProtoVersion, error) {
	response, err := provisioner.transmitter.Send(VersionEndpoint, []byte("---"))

	if err != nil {
		return ProtoVersion{}, err
	}

	var protoVer ProtoVersion
	err = json.Unmarshal(response, &protoVer)

	if err != nil {
		return ProtoVersion{}, err
	}

	return protoVer, nil
}

func (provisioner *Provisioner) EstablishSession() error {
	var response []byte

	for {
		request, err := provisioner.security.SecuritySession(response)
		if err != nil {
			return err
		}

		if request == nil {
			return nil
		}

		response, err = provisioner.transmitter.Send(SessionEndpoint, request)
		if err != nil {
			return err
		}
	}
}

func (provisioner *Provisioner) WiFiScan() ([]WiFiFromScan, error) {
	msg, err := StartScanRequest(provisioner.security, true, false, 5, 120)
	if err != nil {
		return nil, err
	}

	resp, err := provisioner.transmitter.Send(ScanEndpoint, msg)

	if err != nil {
		return nil, err
	}

	err = StartScanResponse(provisioner.security, resp)
	if err != nil {
		return nil, err
	}

	msg, err = ScanStatusRequest(provisioner.security)
	if err != nil {
		return nil, err
	}

	resp, err = provisioner.transmitter.Send(ScanEndpoint, msg)
	if err != nil {
		return nil, err
	}

	scanStatusResult, err := ScanStatusResponse(provisioner.security, resp)
	if err != nil {
		return nil, err
	}

	wiFis := make([]WiFiFromScan, 0)

	if scanStatusResult.Count != 0 {
		var index uint32 = 0
		remaining := scanStatusResult.Count

		for remaining > 0 {
			count := remaining
			if remaining > 100 {
				count = 100
			}

			msg, err = ScanResultRequest(provisioner.security, index, count)
			if err != nil {
				return nil, err
			}

			resp, err = provisioner.transmitter.Send(ScanEndpoint, msg)
			if err != nil {
				return nil, err
			}

			scanResultResponse, err := ScanResultResponse(provisioner.security, resp)
			if err != nil {
				return nil, err
			}

			for _, wiFi := range scanResultResponse {
				wiFis = append(wiFis, wiFi)
			}

			remaining -= count
			index += count
		}
	}

	return wiFis, nil
}
