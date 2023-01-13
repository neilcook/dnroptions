package options

import (
	"fmt"
	"github.com/miekg/dns"
	_ "gopkg.in/dealancer/validate.v2"
	"net"
	"strings"
)

// RAOption is for V6 addresses
type RAOption struct {
	ServicePriority uint16   `yaml:"svc_prio"`
	Lifetime        uint32   `yaml:"lifetime"`
	ADN             string   `yaml:"adn"`
	Addresses       []net.IP `yaml:"addresses"`
	ServiceParams   string   `yaml:"svc_params" validate:"empty=false"`
}

type RAOptions struct {
	RAOptions []*RAOption `yaml:"ra_options"`
}

func (d *RAOption) Validate() error {
	if d.ADN == "" {
		return fmt.Errorf("ADN is mandatory")
	}
	if !strings.HasSuffix(d.ADN, ".") {
		d.ADN = d.ADN + "."
	}
	// Check ADN is a valid DNS name
	if !dns.IsFqdn(d.ADN) {
		return fmt.Errorf("ADN %s is not a valid DNS name", d.ADN)
	}
	for _, a := range d.Addresses {
		if a.To4() != nil {
			return fmt.Errorf("Address %v is not an IPv6 address", a)
		}
	}
	// Validate ServiceParams according to Section 2.1 of [I-D.ietf-dnsop-svcb-https] if there are addresses
	if len(d.Addresses) > 0 {
		err := ValidateSvcParams(d.ServiceParams)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *RAOptions) Validate() error {
	for _, s := range d.RAOptions {
		err := s.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *RAOption) Serialize() ([]byte, error) {
	var optBuf []byte
	optBuf = append(optBuf, HostToNetShort(d.ServicePriority)...)
	optBuf = append(optBuf, HostToNetLong(d.Lifetime)...)
	adnBuf := make([]byte, len(d.ADN)+1)
	_, err := dns.PackDomainName(d.ADN, adnBuf, 0, nil, false)
	if err != nil {
		return nil, err
	}
	optBuf = append(optBuf, HostToNetShort(uint16(len(adnBuf)))...)
	optBuf = append(optBuf, adnBuf...)
	var addrLen = 0
	addrLen = 16 * len(d.Addresses)
	optBuf = append(optBuf, HostToNetShort(uint16(addrLen))...)
	for _, a := range d.Addresses {
		optBuf = append(optBuf, a.To16()...)
	}
	optBuf = append(optBuf, HostToNetShort(uint16(len(d.ServiceParams)))...)
	optBuf = append(optBuf, []byte(d.ServiceParams)...)
	padBytes := (len(optBuf) + 2) % 8
	for i := 0; i < padBytes; i++ {
		optBuf = append(optBuf, uint8(0))
	}
	return optBuf, nil
}

func (d *RAOptions) Serialize() ([]byte, error) {
	var optBuf []byte
	for _, s := range d.RAOptions {
		b, err := s.Serialize()
		if err != nil {
			return nil, err
		}
		optBuf = append(optBuf, b...)
	}
	return optBuf, nil
}
