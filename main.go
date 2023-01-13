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

func main() {
	var configFile = flag.String("config", "config.yaml", "YAML config file")
	flag.Parse()

	options, err := LoadOptions(*configFile)
	if err != nil {
		log.Fatalf("FATAL: Could not load -config %s : %v", *configFile, err)
	}
	encodedOpts, err := options.DHCPOptions.Serialize()
	if err != nil {
		log.Fatalf("FATAL: Could not serialize options: %v", err)
	}
	hexOpts := make([]byte, hex.EncodedLen(len(encodedOpts)))
	hex.Encode(hexOpts, encodedOpts)
	if options.V6 {
		fmt.Printf("DHCPV6=%s\n", hexOpts)
	} else {
		fmt.Printf("DHCPV4=%s\n", hexOpts)
	}
	encodedOpts, err = options.RAOptions.Serialize()
	if err != nil {
		log.Fatalf("FATAL: Could not serialize options: %v", err)
	}
	hexOpts = make([]byte, hex.EncodedLen(len(encodedOpts)))
	hex.Encode(hexOpts, encodedOpts)
	fmt.Printf("IPV6RA=%s\n", hexOpts)
}
