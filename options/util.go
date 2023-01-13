package options

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func HostToNetShort(i uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}

func HostToNetLong(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func ValidateSvcParams(s string) error {
	s1 := strings.Fields(s)
	foundAlpn := false
	for _, t := range s1 {
		t1 := strings.Split(t, "=")
		for _, sp := range t1 {
			if sp == "alpn" {
				foundAlpn = true
			}
		}
	}
	if foundAlpn {
		return nil
	} else {
		return fmt.Errorf("Service Params (%s) do not contain alpn", s)
	}
}
