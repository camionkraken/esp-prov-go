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
