package softap

import (
	"esp-prov-go/core"
	"esp-prov-go/core/security"
)

func NewSoftapProvisioner(hostname string, security security.Security) *core.Provisioner {
	actualHostname := DefaultHostname

	if hostname != "" {
		actualHostname = hostname
	}

	return core.NewProvisioner(NewHttpTransmitter(actualHostname), security)
}
