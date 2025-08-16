package softap

import (
	"esp-prov-go/core"
	"esp-prov-go/core/security"
)

func NewSoftapProvisioner(hostname string, security security.Security) (*core.Provisioner, error) {
	actualHostname := DefaultHostname

	if hostname != "" {
		actualHostname = hostname
	}

	transmitter, err := NewHttpTransmitter(actualHostname)

	if err != nil {
		return nil, err
	}

	return core.NewProvisioner(transmitter, security), nil
}
