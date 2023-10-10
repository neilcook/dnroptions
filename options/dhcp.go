package options

import (
	"fmt"
	"github.com/miekg/dns"
	_ "gopkg.in/dealancer/validate.v2"
	"net"
	"strings"
)

// DHCPOption is for V6 or V4 DNR Options
type DHCPOption struct {
	V6              bool     `yaml:"-"`
	ServicePriority uint16   `yaml:"svc_prio"`
	ADN             string   `yaml:"adn"`
	Addresses       []net.IP `yaml:"addresses"`
	ServiceParams   string   `yaml:"svc_params"`
}

type DHCPOptions struct {
	DHCPOptions []*DHCPOption `yaml:"dhcp_options""`
	V6          bool          `yaml:"v6"`
}

func (d *DHCPOption) Validate() error {
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
		if d.V6 && a.To4() != nil {
			return fmt.Errorf("Address %v is not an IPv6 address", a)
		} else if !d.V6 && a.To4() == nil {
			return fmt.Errorf("Address %v is not an IPv4 address", a)
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

func (d *DHCPOptions) Validate() error {
	for _, s := range d.DHCPOptions {
		if d.V6 {
			s.V6 = true
		}
		err := s.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DHCPOption) Serialize() ([]byte, error) {
	var optBuf []byte
	// DHCPv4 has a 2 byte length parameter
	if !d.V6 {
		optBuf = append(optBuf, HostToNetShort(0)...)
	}
	optBuf = append(optBuf, HostToNetShort(d.ServicePriority)...)
	adnBuf := make([]byte, len(d.ADN)+1)
	_, err := dns.PackDomainName(d.ADN, adnBuf, 0, nil, false)
	if err != nil {
		return nil, err
	}
	if d.V6 {
		optBuf = append(optBuf, HostToNetShort(uint16(len(adnBuf)))...)
	} else {
		optBuf = append(optBuf, uint8(len(adnBuf)))
	}
	optBuf = append(optBuf, adnBuf...)
	var addrLen = 0
	if d.V6 {
		addrLen = 16 * len(d.Addresses)
	} else {
		addrLen = 4 * len(d.Addresses)
	}
	if d.V6 {
		optBuf = append(optBuf, HostToNetShort(uint16(addrLen))...)
	} else {
		optBuf = append(optBuf, HostToNetByte(uint8(addrLen))...)
	}
	// IPv4 is only 8 bits
	for _, a := range d.Addresses {
		if d.V6 {
			optBuf = append(optBuf, a.To16()...)
		} else {
			optBuf = append(optBuf, a.To4()...)
		}
	}
	optBuf = append(optBuf, []byte(d.ServiceParams)...)
	// Update the length parameter for V4
	if !d.V6 {
		copy(optBuf[0:2], HostToNetShort(uint16(len(optBuf[2:]))))
	}
	return optBuf, nil
}

func (d *DHCPOptions) Serialize() ([]byte, error) {
	var optBuf []byte
	for _, s := range d.DHCPOptions {
		b, err := s.Serialize()
		if err != nil {
			return nil, err
		}
		optBuf = append(optBuf, b...)
	}
	return optBuf, nil
}
