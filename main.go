package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/neilcook/dnroptions/options"
	"log"
	"os"
)

type DNROptions struct {
	options.DHCPOptions `yaml:",inline"`
	options.RAOptions   `yaml:",inline"`
}

func LoadFromBytes(yamlContents []byte) (*DNROptions, error) {
	d := DNROptions{}
	err := yaml.Unmarshal(yamlContents, &d)
	if err != nil {
		return nil, err
	}
	err = d.DHCPOptions.Validate()
	if err != nil {
		return nil, err
	}
	err = d.RAOptions.Validate()
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func LoadOptions(fpath string) (*DNROptions, error) {
	yamlFile, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	d, err := LoadFromBytes(yamlFile)
	return d, err
}

func hexEncode(b []byte, colon bool) []byte {
	hexOpts := make([]byte, 0, 3*len(b))
	x := hexOpts[1*len(b) : 3*len(b)]
	hex.Encode(x, b)
	if colon {
		for i := 0; i < len(x); i += 2 {
			hexOpts = append(hexOpts, x[i], x[i+1], ':')
		}
		hexOpts = hexOpts[:len(hexOpts)-1]
	} else {
		hexOpts = x
	}
	return hexOpts
}

func main() {
	var configFile = flag.String("config", "config.yaml", "YAML config file")
	var colonHex = flag.Bool("hexcolons", true, "Use colons to separate hex octets in output")
	flag.Parse()

	options, err := LoadOptions(*configFile)
	if err != nil {
		log.Fatalf("FATAL: Could not load -config %s : %v", *configFile, err)
	}
	encodedOpts, err := options.DHCPOptions.Serialize()
	if err != nil {
		log.Fatalf("FATAL: Could not serialize options: %v", err)
	}
	hexOpts := hexEncode(encodedOpts, *colonHex)
	if options.V6 {
		fmt.Printf("DHCPV6=%s\n", hexOpts[:len(hexOpts)])
	} else {
		fmt.Printf("DHCPV4=%s\n", hexOpts)
	}
	encodedOpts, err = options.RAOptions.Serialize()
	if err != nil {
		log.Fatalf("FATAL: Could not serialize options: %v", err)
	}
	hexOpts = hexEncode(encodedOpts, *colonHex)
	fmt.Printf("IPV6RA=%s\n", hexOpts)
}
