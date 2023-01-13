# dnroptions
Go Tool to create DHCP and RA Options for IETF DNR Standard

Overview
---

The IETF DNR Standard specifies new DHCP and IPv6 Router Advertisement options
to discover encrypted DNS resolvers. However as an author of this standard, when I tried to
add these options to a DHCP server, I realised these options are in a binary format that it 
could be difficult for many users to easily create in a form that can be put into the 
configuration of a DHCP server.

The dnroptions tool enables you to put the required parameters into a simple YAML configuration file,
and it then outputs Hex-encoded options which can be easily inserted into the configuration of DHCP
servers such as ISC Kea or similar software.

The tool supports DHCPv6, DHCPv4 and RAv6 DNR options. It supports both single and muliple options. 
For multiple options, it will output a single hex-encoded string with all the options concatenated 
together; If you do not want this, then simply run the tool multiple times over different configuration files.

The tool only output the option data itself, not the Option Code or the Option Length - these would
be emitted by the server itself. For the RAv6 option, it does pad the option with zeros to an 8-octet 
boundary as if the the 2-octet type and length field were added (i.e. once you add those 2 octets to the
front then it is 64-bit aligned).

The config file for DHCP options looks like:
````
dhcp_options:
  - svc_prio: 10
    adn: "foobar.com"
    addresses:
      - "127.0.0.1"
      - "192.243.2.1"
    svc_params: alpn=h2,h3
````
Running `dnroptions --config file_with_the_above_config.yaml` will produce the following output:
````
DHCPV4=0023000a0c06666f6f62617203636f6d0000087f000001c0f30201616c706e3d68322c6833
IPV6RA=
````
The default is for IPv4 options; it will complain if it sees IPv6 addresses. To create v6 options, use
the following:
````
v6: true
dhcp_options:
  - svc_prio: 10
    adn: "foobar.com"
    addresses:
      - "fc0e::"
      - "ae31::"
    svc_params: alpn=h2,h3
````
Running `dnroptions --config file_with_the_above_config.yaml` will produce the following output:
````
DHCPV6=000a000c06666f6f62617203636f6d000020fc0e0000000000000000000000000000ae310000000000000000000000000000616c706e3d68322c6833
IPV6RA=
````
To produce RAv6 options, use the following config:
````
ra_options:
  - svc_prio: 10
    lifetime: 123445
    adn: "foobar.com"
    addresses:
      - "fc0e::"
      - "ae31::"
    svc_params: alpn=h2,h3
````
Running `dnroptions --config file_with_the_above_config.yaml` will produce the following output:
````
DHCPV4=
IPV6RA=000a0001e235000c06666f6f62617203636f6d000020fc0e0000000000000000000000000000ae310000000000000000000000000000000a616c706e3d68322c683300000000
````
You can specify multiple options at once, and the tool also supports specifying only the ADN. For example:
````
dhcp_options:
  - svc_prio: 10
    adn: "foobar.com"
    addresses:
      - "127.0.0.1"
      - "192.243.2.1"
    svc_params: alpn=h2,h3
  - adn: "wibble.com"
    addresses:
      - "128.243.2.1"
      - "8.8.8.8"
    svc_params: alpn=h2 port=1234
  - adn: "boo.com"
````
Will produce:
````
DHCPV4=0023000a0c06666f6f62617203636f6d0000087f000001c0f30201616c706e3d68322c6833002a00000c06776962626c6503636f6d00000880f3020108080808616c706e3d683220706f72743d31323334000e00000903626f6f03636f6d000000
IPV6RA=
````
