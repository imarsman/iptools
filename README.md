# iptools

This project is a learning experience for ipv4 subnetting and ipv6 addresses. Mostly I am interested in learning more
about ipv4 subnets and splitting up subnets into equal ranges and defining things like the network ID and broadcast
address for subnets. This is a way for me to better establish my knowledge of this subject.

Given that this code represents me learning I do not recommend using it for production purposes, as the standards being
implemented are very strict for that. I may easily have made errors that I have yet to catch and fix.

This utility can have completion enabled by typing `COMP_INSTALL=1 iptools`.

This utility currently does three things with ipv4. It splits a subnet into networks and into networks by a differing subnet size,
it splits a subnet into a set of ranges for its networks, and it gives summary information for a subnet.

For ipv6 a random ipv6 IP can be generated (or manually entered as a parameter) and described. Currently global unicast,
link local, and private addresses are handled along with multicast, interface local multicast, and link local multicast.
The full range of address types will at tome point be supported where it makes sense. 

Especially with ipv6 addresses contain a great amount of information that relies on external software to keep track of.
It is possible to produce valid ipv6 addresses but the random numbers used to generate the addresses are in no way
meaningfully tied to any system that might use them even if they are "valid". One possible exception would be Link local
addresses which rely on a good random number generating algorithm for the interface ID.

The IP package used, `net/netip`, is in Go 1.18. It is slightly different in API compared to the
[netaddr](https://github.com/inetaf/netaddr) package that came first. Sadly, the `netip` package loses the IPRange
struct, so I have added needed functionality here in the ipv4 subnet package.

I would be surprised if there were not errors due to coding or to a lack of understanding of ipv4 and ipv6.

### Subnet divisions split default

```
$ iptools subnetip4 ranges -ip 10.32.0.0 -bits 16 -pretty

     Category              Value
------------------- --------------------
 Subnet              10.32.0.0/16
 Subnet IP           10.32.0.0
 Broadcast Address   10.32.255.255
 Subnet Mask         255.255.0.0
 Wildcard Mask       0.0.255.255
 Networks            1
 Network Hosts       %!d(string=65,536)

   Start          End
----------- ---------------
 10.32.0.0   10.32.255.255
 ```

### Subnet divisions split into non-default sized networks

```
$ iptools subnetip4 divide -ip 10.32.0.0 -bits 16 -secondary-bits 18 -pretty

              Category                    Value
------------------------------------ ---------------
 Subnet                               10.32.0.0/16
 Subnet IP                            10.32.0.0
 Broadcast Address                    10.32.255.255
 Subnet Mask                          255.255.0.0
 Wildcard Mask                        0.0.255.255
 Secondary Subnet                     10.32.0.0/18
 Secondary Subnet IP                  10.32.0.0
 Secondary Subnet Broadcast Address   10.32.63.255
 Secondary Subnet Mask                255.255.192.0
 Secondary Subnet Wildcard Mask       0.0.63.255
 Networks                             1
 Secondary Networks                   4
 Effective Networks                   4
 Network Hosts                        65536
 Sub Network Hosts                    16384

    Subnets
----------------
 10.32.0.0/18
 10.32.64.0/18
 10.32.128.0/18
 10.32.192.0/18
 ```

### Subnet ranges in default size

```
$ iptools subnetip4 ranges -ip 99.236.32.0 -bits 16 -pretty

     Category            Value
------------------- ----------------
 Subnet              99.236.0.0/16
 Subnet IP           99.236.0.0
 Broadcast Address   99.236.255.255
 Subnet Mask         255.255.0.0
 Networks            1
 Network Hosts       65536

   Start           End
------------ ----------------
 99.236.0.0   99.236.255.255
 ```
 
### Subnet ranges split into non-default size

```
$ iptools subnetip4 ranges -ip 99.236.32.0 -bits 16 -secondary-bits 18 -pretty

              Category                    Value
------------------------------------ ----------------
 Subnet                               99.236.0.0/16
 Subnet IP                            99.236.0.0
 Broadcast Address                    99.236.255.255
 Subnet Mask                          255.255.0.0
 Wildcard Mask                        0.0.255.255
 Secondary Subnet                     99.236.0.0/18
 Secondary Subnet IP                  99.236.0.0
 Secondary Subnet Broadcast Address   99.236.63.255
 Secondary Subnet Mask                255.255.192.0
 Secondary Subnet Wildcard Mask       0.0.63.255
 Networks                             1
 Secondary Networks                   4
 Effective Networks                   4
 Network Hosts                        65,536
 Secondary Network Hosts              16,384

    Start            End
-------------- ----------------
 99.236.0.0     99.236.63.255
 99.236.64.0    99.236.127.255
 99.236.128.0   99.236.191.255
 99.236.192.0   99.236.255.255
```


### Subnet details

```
$ iptools subnetip4 describe -ip 10.32.0.0 -bits 23

         Category                          Value
-------------------------- -------------------------------------
 Subnet                     10.32.0.0/23
 Subnet IP                  10.32.0.0
 Broadcast Address          10.32.1.255
 Broadcast Address Hex ID   0xA2001FF
 Subnet Mask                255.255.254.0
 Wildcard Mask              0.0.1.255
 IP Class                   A
 IP Type                    Private
 Binary Subnet Mask         00001010.00100000.00000000.00000000
 Binary ID                  00001010001000000000000000000000
 in-addr.arpa               0.0.32.10.in-addr.arpa
 Networks                   128
 Network Hosts              512
```

#### Subnet with secondary subnet details

```
$ iptools subnetip4 describe -ip 10.32.0.0 -bits 16 -secondary-bits 18

              Category                               Value
------------------------------------ -------------------------------------
 Subnet                               10.32.0.0/16
 Subnet IP                            10.32.0.0
 Broadcast Address                    10.32.255.255
 Broadcast Address Hex ID             0xA20FFFF
 Subnet Mask                          255.255.0.0
 Wildcard Mask                        0.0.255.255
 IP Class                             A
 IP Type                              Private
 Binary Subnet Mask                   00001010.00100000.00000000.00000000
 Binary ID                            00001010001000000000000000000000
 in-addr.arpa                         0.0.32.10.in-addr.arpa
 Secondary Subnet                     10.32.0.0/18
 Secondary Subnet IP                  10.32.0.0
 Secondary Subnet Broadcast Address   10.32.63.255
 Secondary Subnet Mask                255.255.192.0
 Secondary Subnet Wildcard Mask       0.0.63.255
 Networks                             1
 Secondary Networks                   4
 Effective Networks                   4
 Network Hosts                        65,536
 Secondary Network Hosts              16,384
```

### IPV6 Global unicast address

Parse an ip with prefix
```
$ iptools ip6 describe -ip 2001:0db8:85a3:0000:0000:8a2e:0370:7334/64

         Category                                            Value
-------------------------- --------------------------------------------------------------------------
-------------------------- --------------------------------------------------------------------------
 IP Type                    Global unicast
 Type Prefix                2000::/3
 IP                         2001:db8:85a3::8a2e:370:7334
 Solicited node multicast   ff02::1:ff70:7334
 Prefix                     2001:db8:85a3::/64
 Routing Prefix             2001:0db8:85a3::/48
 Subnet ID                  0000
 Subnets                    65,536
 Global ID                  01:0db8:85a3
 Interface ID               0000:8a2e:0370:7334
 Addresses                  18,446,744,073,709,551,616
 Link                       http://[2001:db8:85a3::8a2e:370:7334]/
 ip6.arpa                   4.3.3.7.0.7.3.0.e.2.a.8.0.0.0.0.0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa
 Subnet first address       2001:0db8:85a3:0000:0000:0000:0000:0000
 Subnet last address        2001:0db8:85a3:0000:ffff:ffff:ffff:ffff
 1st address field binary   0010000000000001
```

Parse an IP and supply bits separately
```
$ iptools ip6 describe -ip 2001:0db8:85a3:0000:0000:8a2e:0370:7334 -bits 64
         Category                                            Value
-------------------------- --------------------------------------------------------------------------
 IP Type                    Global unicast
 Type Prefix                2000::/3
 IP                         2001:db8:85a3::8a2e:370:7334
 Solicited node multicast   ff02::1:ff70:7334
 Prefix                     2001:db8:85a3::/64
 Routing Prefix             2001:0db8:85a3::/48
 Subnet ID                  0000
 Subnets                    65,536
 Global ID                  01:0db8:85a3
 Interface ID               0000:8a2e:0370:7334
 Addresses                  18,446,744,073,709,551,616
 Link                       http://[2001:db8:85a3::8a2e:370:7334]/
 ip6.arpa                   4.3.3.7.0.7.3.0.e.2.a.8.0.0.0.0.0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa
 Subnet first address       2001:0db8:85a3:0000:0000:0000:0000:0000
 Subnet last address        2001:0db8:85a3:0000:ffff:ffff:ffff:ffff
 1st address field binary   0010000000000001
```

```
$ iptools ip6 describe -random -type global-unicast
         Category                                            Value
-------------------------- --------------------------------------------------------------------------
 IP Type                    Global unicast
 Type Prefix                2000::/3
 IP                         2301:db8:cafe:bd9e:72f4:92ff:fe5b:a76
 Solicited node multicast   ff02::1:ff5b:a76
 Routing Prefix             2301:0db8:cafe::/48
 Subnet ID                  bd9e
 Subnets                    65,536
 Global ID                  301:0db8:cafe
 Interface ID               72f4:92ff:fe5b:0a76
 Addresses                  18,446,744,073,709,551,616
 Link                       http://[2301:db8:cafe:bd9e:72f4:92ff:fe5b:a76]/
 ip6.arpa                   6.7.a.0.b.5.e.f.f.f.2.9.4.f.2.7.e.9.d.b.e.f.a.c.8.b.d.0.1.0.3.2.ip6.arpa
 Subnet first address       2301:0db8:cafe:bd9e:0000:0000:0000:0000
 Subnet last address        2301:0db8:cafe:bd9e:ffff:ffff:ffff:ffff
 1st address field binary   0010001100000001
```

### IPV6 Link local address
```
$ iptools ip6 describe -random -type link-local
         Category                            Value
-------------------------- -----------------------------------------
 IP Type                    Link local unicast
 Type Prefix                fe80::/10
 IP                         fe80::3aff:afff:feab:9edc
 Solicited node multicast   ff02::1:ffab:9edc
 Subnet ID                  0000
 Subnets                    65,536
 Interface ID               3aff:afff:feab:9edc
 Addresses                  18,446,744,073,709,551,616
 Default Gateway            fe80::1
 Subnet first address       fe80:0000:0000:0000:0000:0000:0000:0000
 Subnet last address        fe80:0000:0000:0000:ffff:ffff:ffff:ffff
 1st address field binary   1111111010000000
 ```

### IPV6 private address
```
$ iptools ip6 describe -random -type private
         Category                            Value
-------------------------- -----------------------------------------
 IP Type                    Private
 Type Prefix                fc00::/7
 IP                         fdfc:3ae0:b79f:413c:b2f5:7dff:fe8a:fcd
 Solicited node multicast   ff02::1:ff8a:fcd
 Subnet ID                  413c
 Subnets                    65,536
 Global ID                  fc:3ae0:b79f
 Interface ID               b2f5:7dff:fe8a:0fcd
 Addresses                  18,446,744,073,709,551,616
 Subnet first address       fdfc:3ae0:b79f:413c:0000:0000:0000:0000
 Subnet last address        fdfc:3ae0:b79f:413c:ffff:ffff:ffff:ffff
 1st address field binary   1111110111111100
```

### IPV6 multicast
```
$ iptools ip6 describe -bits 64 -random -type multicast
          Category                           Value
---------------------------- --------------------------------------
 IP Type                      Multicast
 Type Prefix                  ff00::/8
 IP                           ff13:0:8be8:1642:be9a:d38f:4fc3:aecb
 Network Prefix               8be8:1642:be9a:d38f
 Group ID                     4fc3:aecb
 Groups                       4,294,967,296
 first address field binary   1111111100010011
```

### IPV6 Interface local multicast
```
$ iptools ip6 describe -bits 64 -random -type interface-local-multicast
          Category                           Value
---------------------------- --------------------------------------
 IP Type                      Interface local multicast
 Type Prefix                  ff00::/8
 IP                           ff11:0:d1fb:979f:5ff6:56eb:4a3d:a793
 Network Prefix               d1fb:979f:5ff6:56eb
 Group ID                     4a3d:a793
 Groups                       4,294,967,296
 first address field binary   1111111100010001
```

### IPV6 Link local multicast
```
$ iptools ip6 describe -bits 64 -random -type link-local-multicast
          Category                           Value
---------------------------- --------------------------------------
 IP Type                      Link local muticast
 Type Prefix                  ff00::/8
 IP                           ff12:0:c957:512a:f0fe:bbed:3d57:9c0b
 Network Prefix               c957:512a:f0fe:bbed
 Group ID                     3d57:9c0b
 Groups                       4,294,967,296
 first address field binary   1111111100010010
```

### Generate random IPs
```
$ iptools subnetip6 random-ips -number 10 -type unique-local
fd00:0fc3:d81c:ebbd:211d:58ff:fe21:054b
fd00:5b7e:2277:4ace:2171:51ff:fe8b:b8bb
fd00:07a9:cf42:f39f:bf25:b9ff:fe57:bd24
fd00:2dc6:322e:e46c:023f:b1ff:fea8:9650
fd00:ee0f:dde0:2d8b:f841:c0ff:fe3c:8665
fd00:5d4a:4b6c:ba79:8c74:66ff:fee1:4115
fd00:3340:3705:3a15:36bd:eaff:fee1:6e9d
fd00:0092:8f02:e97a:8eef:b4ff:feeb:28bb
fd00:a954:b6c1:3e94:c070:e3ff:fe9c:42ea
fd00:6df2:ea4f:5d05:52cb:feff:fe81:e3b1
```

## Utilities

### Lookup of IPs by domain

```
$ iptools utilities lookup-domains -domains cisco.com ibm.com
 Type           Address
------ -------------------------
        cisco.com
 ipv4   72.163.4.185
 ipv6   2001:420:1101:1::185
        ibm.com
 ipv4   104.67.113.240
 ipv6   2607:f798:d04:283::3831
 ipv6   2607:f798:d04:289::3831
```

### Top level help

```
$ iptools -h
iptools
-------
Commit:  b3fd0d9
Date:    2022-09-21T00:55:35Z
Tag:     v0.1.21
OS:      darwin
ARCH:    arm64

Usage: iptools <command> [<args>]

Options:
  --help, -h             display this help and exit
  --version              display version and exit

Commands:
  subnetip4              Get networks for subnet
  ip6                    Get IP6 address information
```

### Help for subnet options

```
$ iptools subnetip4 -h
iptools
-------
Commit:  b3fd0d9
Date:    2022-09-21T00:55:35Z
Tag:     v0.1.21
OS:      darwin
ARCH:    arm64

Usage: iptools subnetip4 <command> [<args>]
  --help, -h             display this help and exit
  --version              display version and exit

Commands:
  ranges                 divide a subnet into ranges
  divide                 divide a subnet into smaller subnets
  describe               describe a subnet

$ iptools ip6 describe -h
iptools
-------
Commit:  954eb63
Date:    2022-09-22T03:02:28Z
Tag:     v0.1.21
OS:      darwin
ARCH:    arm64

Usage: iptools ip6 describe [--ip IP] [--random] [--bits BITS] [--type TYPE]

Options:
  --ip IP, -i IP         IP address
  --random, -r           generate random IP
  --bits BITS, -b BITS   subnet bits
  --type TYPE, -t TYPE   global-unicast, link-local, unique-local
  --help, -h             display this help and exit
  --version              display version and exit
```

### Lines of code

```
$ gocloc pkg cmd README.md
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                              10            367            265           1976
Markdown                         1             53              0            320
-------------------------------------------------------------------------------
TOTAL                           11            420            265           2296
-------------------------------------------------------------------------------
```
