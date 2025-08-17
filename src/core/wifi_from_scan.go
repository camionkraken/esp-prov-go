package core

import (
	"errors"
	"esp-prov-go/core/proto/protogen"
	"fmt"
)

type WiFiFromScan struct {
	Ssid    string
	BSsid   string
	Channel *uint32
	Rssi    *int32
	Auth    *WiFiAuth
}

func FromWiFiScanResult(wiFiScanResult *protogen.WiFiScanResult) (WiFiFromScan, error) {
	if wiFiScanResult == nil {
		return WiFiFromScan{}, errors.New("wiFiScanResult is nil")
	}

	auth, err := fromWiFiAuthMode(&wiFiScanResult.Auth)

	if err != nil {
		return WiFiFromScan{}, err
	}
	
	return WiFiFromScan{
		Ssid: string(wiFiScanResult.Ssid),
		BSsid: string(wiFiScanResult.Bssid),
		Channel: &wiFiScanResult.Channel,
		Rssi: &wiFiScanResult.Rssi,
		Auth: &auth,
	}, nil
}

type WiFiAuth int

const (
	WiFiAuthOpen WiFiAuth = iota
	WiFiAuthWep
	WiFiAuthWpaPsk
	WiFiAuthWpa2Psk
	WiFiAuthWpaWpa2Psk
	WiFiAuthWpa2Enterprise
	WiFiAuthWpa3Psk
	WiFiAuthWpa2Wpa3Psk
)

func fromWiFiAuthMode(wiFiAuthMode *protogen.WifiAuthMode) (WiFiAuth, error) {
	switch *wiFiAuthMode {
	case protogen.WifiAuthMode_Open:
		return WiFiAuthOpen, nil
	case protogen.WifiAuthMode_WEP:
		return WiFiAuthWep, nil
	case protogen.WifiAuthMode_WPA_PSK:
		return WiFiAuthWpaPsk, nil
	case protogen.WifiAuthMode_WPA2_PSK:
		return WiFiAuthWpa2Psk, nil
	case protogen.WifiAuthMode_WPA_WPA2_PSK:
		return WiFiAuthWpaWpa2Psk, nil
	case protogen.WifiAuthMode_WPA2_ENTERPRISE:
		return WiFiAuthWpa2Enterprise, nil
	case protogen.WifiAuthMode_WPA3_PSK:
		return WiFiAuthWpa3Psk, nil
	case protogen.WifiAuthMode_WPA2_WPA3_PSK:
		return WiFiAuthWpa2Wpa3Psk, nil
	default:
		return WiFiAuth(0), errors.New(fmt.Sprintf("invalid WiFiAuthMode: %d", wiFiAuthMode))
	}
}