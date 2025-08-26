package core

import (
	"encoding/json"
	"esp-prov-go/core/security"
	"fmt"
	"time"
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

func (provisioner *Provisioner) GetWiFiStatus() (ProvStatusResult, error) {
	msg, err := GetStatusRequest(provisioner.security)
	if err != nil {
		return ProvStatusUnknown, err
	}

	response, err := provisioner.transmitter.Send(ProvConfigEndpoint, msg)
	if err != nil {
		return ProvStatusUnknown, err
	}

	status, err := GetStatusResponse(provisioner.security, response)
	if err != nil {
		return ProvStatusUnknown, err
	}

	return status, nil
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

func (provisioner *Provisioner) ConnectToWiFiNetwork(ssid string, passphrase string) (WiFiConnectionResult, error) {
	msg, err := SetConfigRequest(provisioner.security, ssid, passphrase)
	if err != nil {
		return WiFiSetConfigurationFailed, err
	}

	resp, err := provisioner.transmitter.Send(ProvConfigEndpoint, msg)
	if err != nil {
		return WiFiSetConfigurationFailed, err
	}

	err = SetConfigResponse(provisioner.security, resp)
	if err != nil {
		return WiFiSetConfigurationFailed, err
	}

	msg, err = ApplyConfigRequest(provisioner.security)
	if err != nil {
		return WiFiApplyConfigurationFailed, err
	}

	resp, err = provisioner.transmitter.Send(ProvConfigEndpoint, msg)
	if err != nil {
		return WiFiApplyConfigurationFailed, err
	}

	if err = ApplyConfigResponse(provisioner.security, resp); err != nil {
		return WiFiApplyConfigurationFailed, err
	}

	for i := range ConnectionCheckRetries {
		status, err := provisioner.GetWiFiStatus()
		if err != nil {
			fmt.Printf("Error getting WiFi status: %s. Try %d\n", err, i)
			continue
		}

		switch status {
		case ProvStatusConnecting:
			time.Sleep(ConnectionCheckIntervalMs * time.Millisecond)
			continue
		case ProvStatusConnected:
			return WiFiConnected, nil
		case ProvStatusNetworkNotFound:
			return WiFiNetworkNotFound, nil
		case ProvStatusAuthError:
			return WiFiAuthError, nil
		case ProvStatusUnknown:
		case ProvStatusFailed:
		case ProvStatusDisconnected:
			return WiFiFailed, nil
		}
	}

	return WiFiRetriesExceeded, nil
}
