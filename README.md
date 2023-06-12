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


     Category            Value
------------------- ---------------
 Subnet              10.32.0.0/16
 Subnet IP           10.32.0.0
 Broadcast Address   10.32.255.255
 Subnet Mask         255.255.0.0
 Wildcard Mask       0.0.255.255
 Networks            1
 Network Hosts       65,536

   Start          End
----------- ---------------
 10.32.0.0   10.32.255.255
 ```

### Subnet divisions split into non-default sized networks

```
 iptools subnetip4 divide -ip 10.32.0.0 -bits 16 -secondary-bits 18 -pretty


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
 Netorks                              1
 Secondary Networks                   4
 Effective Networks                   4
 Network Hosts                        65,536
 Secondary Network Hosts              16,384

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
 IP Type             Public
 Subnet              99.236.0.0/16
 Subnet IP           99.236.0.0
 Broadcast Address   99.236.255.255
 Subnet Mask         255.255.0.0
 Wildcard Mask       0.0.255.255
 Networks            1
 Network Hosts       65,536

   Start           End
------------ ----------------
 99.236.0.0   99.236.255.255
 ```
 
### Subnet ranges split into non-default size

```
$ iptools subnetip4 ranges -ip 99.236.32.0 -bits 16 -secondary-bits 18 -pretty

              Category                    Value
------------------------------------ ----------------
 IP Type                              Public
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
 IP Type                    Private
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
 IP Type                              Private
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
 IP                         3701:db8:cafe:b8cb:72cf:8aff:fe3a:fa69
 Solicited node multicast   ff02::1:ff3a:fa69
 Prefix                     3701:db8:cafe:b8cb::/64
 Routing Prefix             3701:0db8:cafe::/48
 Subnet ID                  b8cb
 Subnets                    65,536
 Global ID                  701:0db8:cafe
 Interface ID               72cf:8aff:fe3a:fa69
 Addresses                  18,446,744,073,709,551,616
 Link                       http://[3701:db8:cafe:b8cb:72cf:8aff:fe3a:fa69]/
 ip6.arpa                   9.6.a.f.a.3.e.f.f.f.a.8.f.c.2.7.b.c.8.b.e.f.a.c.8.b.d.0.1.0.7.3.ip6.arpa
 Subnet first address       3701:0db8:cafe:b8cb:0000:0000:0000:0000
 Subnet last address        3701:0db8:cafe:b8cb:ffff:ffff:ffff:ffff
 1st address field binary   0011011100000001
 ```

### IPV6 Link local address
```
$ iptools ip6 describe -random -type link-local
         Category                            Value
-------------------------- -----------------------------------------
 IP Type                    Link local unicast
 Type Prefix                fe80::/10
 IP                         fe80::7263:80ff:fe2e:d2ff
 Solicited node multicast   ff02::1:ff2e:d2ff
 Prefix                     fe80::/64
 Subnet ID                  0000
 Subnets                    65,536
 Interface ID               7263:80ff:fe2e:d2ff
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
 IP                         fd14:761a:1f7a:e2e9:1af6:97ff:fe1b:1342
 Solicited node multicast   ff02::1:ff1b:1342
 Prefix                     fd14:761a:1f7a:e2e9::/64
 Subnet ID                  e2e9
 Subnets                    65,536
 Global ID                  14:761a:1f7a
 Interface ID               1af6:97ff:fe1b:1342
 Addresses                  18,446,744,073,709,551,616
 Subnet first address       fd14:761a:1f7a:e2e9:0000:0000:0000:0000
 Subnet last address        fd14:761a:1f7a:e2e9:ffff:ffff:ffff:ffff
 1st address field binary   1111110100010100
```

### IPV6 multicast
```
$ iptools ip6 describe -bits 64 -random -type multicast
          Category                           Value
---------------------------- -------------------------------------
 IP Type                      Multicast
 Type Prefix                  ff00::/8
 IP                           ff14:0:5a1c:9ae:ef85:33ae:737e:6bd6
 Prefix                       ff14:0:5a1c:9ae::/64
 Network Prefix               5a1c:09ae:ef85:33ae
 Group ID                     737e:6bd6
 Groups                       4,294,967,296
 first address field binary   1111111100010100
```

### IPV6 Interface local multicast
```
$ iptools ip6 describe -bits 64 -random -type interface-local-multicast
          Category                           Value
---------------------------- -------------------------------------
 IP Type                      Interface local multicast
 Type Prefix                  ff00::/8
 IP                           ff31:0:8e37:805a:438e:ee6c:3f0d:4e8
 Prefix                       ff31:0:8e37:805a::/64
 Network Prefix               8e37:805a:438e:ee6c
 Group ID                     3f0d:04e8
 Groups                       4,294,967,296
 first address field binary   1111111100110001
```

### IPV6 Link local multicast
```
$ iptools ip6 describe -bits 64 -random -type link-local-multicast
          Category                           Value
---------------------------- -------------------------------------
 IP Type                      Link local muticast
 Type Prefix                  ff00::/8
 IP                           ff22:0:e57f:db8:a927:3bbc:49d4:91b4
 Prefix                       ff22:0:e57f:db8::/64
 Network Prefix               e57f:0db8:a927:3bbc
 Group ID                     49d4:91b4
 Groups                       4,294,967,296
 first address field binary   1111111100100010
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
