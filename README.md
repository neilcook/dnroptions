# dnroptions
Go Tool to create DHCP and RA Options for IETF DNR Standard (compliant with RFC 9463)

# Overview

The IETF DNR Standard specifies new DHCP and IPv6 Router Advertisement options
to discover encrypted DNS resolvers. However as an author of this standard, when I tried to
add these options to a DHCP server, I realised these options are in a binary format that it 
could be difficult for many users to easily create in a form that can be put into the 
configuration of a DHCP server.

Update - since I created this tool, some servers such as ISC Kea now support creating some of these options 
in a nicer way for you, however you still have to add the IP address and SvcParams options yourself in a binary fashion,
which makes me wonder why they bothered, but there you go.

The dnroptions tool enables you to put the required parameters into a simple YAML configuration file,
and it then outputs Hex-encoded options which can be easily inserted into the configuration of DHCP
servers such as ISC Kea or similar software. By default the hex octets will be separated with colons;
to use spaces instead, specify `--hexspaces=true`.

The tool supports DHCPv6, DHCPv4 and RAv6 DNR options. It supports both single and muliple options. 
For multiple options, it will output a single hex-encoded string with all the options concatenated 
together; If you do not want this (it's probably not what you want for V6 options, but would be fine for V4 options), 
then simply run the tool multiple times over different configuration files.

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
DHCPV4=00:22:00:0a:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:08:7f:00:00:01:c0:f3:02:01:00:01:00:06:02:68:32:02:68:33
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
DHCPV6=00:0a:00:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:00:20:fc:0e:00:00:00:00:00:00:00:00:00:00:00:00:00:00:ae:31:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:01:00:06:02:68:32:02:68:33
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
IPV6RA=00:0a:00:01:e2:35:00:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:00:20:fc:0e:00:00:00:00:00:00:00:00:00:00:00:00:00:00:ae:31:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:0a:61:6c:70:6e:3d:68:32:2c:68:33:00:00:00:00
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
DHCPV4=00:22:00:0a:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:08:7f:00:00:01:c0:f3:02:01:00:01:00:06:02:68:32:02:68:33:00:25:00:00:0c:06:77:69:62:62:6c:65:03:63:6f:6d:00:08:80:f3:02:01:08:08:08:08:00:01:00:03:02:68:32:00:03:00:02:04:d2:00:0d:00:00:09:03:62:6f:6f:03:63:6f:6d:00:00
````

# Configuring DHCP Servers

Some common DHCP servers can be configured as follows.

## ISC DHCP Server

To add DNRv4 options to dhcp.conf:
````
option dnrv4 code 162 = string;
option dnrv4 00:22:00:0a:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:08:7f:00:00:01:c0:f3:02:01:00:01:00:06:02:68:32:02:68:33;
````

To add DNRv6 options to dhcp6.conf:
````
option dnrv6 code 144 = string;
option dnrv6 00:0a:00:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:00:20:fc:0e:00:00:00:00:00:00:00:00:00:00:00:00:00:00:ae:31:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:01:00:06:02:68:32:02:68:33;
````

## ISC Kea

As discussed above, although Kea does support DNR directly now, most options still require binary configuration, so you
may still want to add the options using this method directly.

To add DNRv4 options to kea-dhcp4.conf:
````
"Dhcp4": {
    // Other Kea configuration is removed for brevity
    "option-def": [
        {
            "name": "dnrv4",
            "code": 162,
            "type": "binary",
            "space": "dhcp4"
        },
        ...
    ],
    "option-data": [
        {
            "name": "dnrv4",
            "space": "dhcp4",
            "csv-format": false,
            "data": "00:22:00:0a:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:08:7f:00:00:01:c0:f3:02:01:00:01:00:06:02:68:32:02:68:33"
        },
        ...
    ],
    ...
}
````

To add DNRv6 options to kea-dhcp6.conf:
````
"Dhcp6": {
    // Other Kea configuration is removed for brevity
    "option-def": [
        {
            "name": "dnrv6",
            "code": 144,
            "type": "binary",
            "space": "dhcp6"
        },
        ...
    ],
    "option-data": [
        {
            "name": "dnrv6",
            "space": "dhcp6",
            "csv-format": false,
            "data": "00:0a:00:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:00:20:fc:0e:00:00:00:00:00:00:00:00:00:00:00:00:00:00:ae:31:00:00:00:00:00:00:00:00:00:00:00:00:00:00:61:6c:70:6e:3d:68:32:2c:68:33"
        },
        {
            "name": "dnrv6",
            "space": "dhcp6",
            "csv-format": false,
            "data": "00:0a:00:0c:06:66:6f:6f:62:61:72:03:63:6f:6d:00:00:20:fc:0e:00:00:00:00:00:00:00:00:00:00:00:00:00:00:ae:31:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:01:00:06:02:68:32:02:68:33"
        },
        ...
    ],
    ...
}
````
The above example shows multiple DNR DHCPV6 options configured.