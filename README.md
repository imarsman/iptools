# iptools

This project was done as a learning experience for things like IP subnetting. Mostly I was interested in learning more
about ipv4 subnets and splitting up subnets into equal ranges and defining things like the network ID and broadcast
address for subnets. This is a way for me to better establish my knowledge of this subject.

Given that this code represents me learning I do not recommend using it for production purposes, as the standards being
implemented are very strict for that. I may easily have made errors that I have yet to catch and fix.

This utility can have completion enabled by typing `COMP_INSTALL=1 iptools`.

This utility currently does three things with IPV4. It splits a subnet into networks and into networks by a differing subnet size,
it splits a subnet into a set of ranges for its networks, and it gives summary information for a subnet.

For IPV6 a random IPV6 IP can be generated (or manually entered as a parameter) and described. Currently global unicast,
link local, and unique local addresses are handled along with multicast, interface local multicast, and
link local multicast. The full range of address types will at tome point be supported where it makes sense. 

Especially with IPV6 addresses contain a great amount of information that relies on external software to keep track of.
It is possible to produce valid IPV6 addresses but the random numbers used to generate the addresses are in no way
meaningfully tied to any system that might use them even if they are "valid". One possible exception would be Link local
addresses which rely on a good random number generating algorithm for the interface ID.

The IP package used, `net/netip`, is in Go 1.18. It is slightly different in API compared to the
[netaddr](https://github.com/inetaf/netaddr) package that came first. Sadly, the `netip` package loses the IPRange
struct, so I have added needed functionality here in the ipv4 subnet package.

I would be surprised if there were not errors due to coding or to a lack of understanding of IPV4 and IPV6. I will work
to reduce both my ignorance and errors.

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

```
$ iptools ip6 describe -bits 64 -type global-unicast -random
          Category                                             Value
---------------------------- --------------------------------------------------------------------------
 IP Type                      Global unicast
 Type Prefix                  2000::/3
 IP                           2e01:db8:cafe:f715:da59:16ff:fe65:46d8
 Prefix                       2e01:db8:cafe:f715::/64
 Routing Prefix               2e01:0db8:cafe::/48
 Global ID                    e01:0db8:cafe
 Interface ID                 da59:16ff:fe65:46d8
 Subnet ID                    f715
 Link                         http://[2e01:db8:cafe:f715:da59:16ff:fe65:46d8]/
 ip6.arpa                     8.d.6.4.5.6.e.f.f.f.6.1.9.5.a.d.5.1.7.f.e.f.a.c.8.b.d.0.1.0.e.2.ip6.arpa
 Subnet first address         2e01:0db8:cafe:f715:0000:0000:0000:0000
 Subnet last address          2e01:0db8:cafe:f715:ffff:ffff:ffff:ffff
 first address field binary   0010111000000001
```

### IPV6 Link local address
```
$ iptools ip6 describe -bits 64 -type link-local -random
          Category                             Value
---------------------------- -----------------------------------------
 IP Type                      Link local unicast
 Type Prefix                  fe80::/10
 IP                           fe80::1568:49ff:fe79:dbd2
 Prefix                       fe80::/64
 Interface ID                 1568:49ff:fe79:dbd2
 Subnet ID                    0000
 Default Gateway              fe80::1
 Link                         http://[fe80::1568:49ff:fe79:dbd2]/
 Subnet first address         fe80:0000:0000:0000:0000:0000:0000:0000
 Subnet last address          fe80:0000:0000:0000:ffff:ffff:ffff:ffff
 first address field binary   1111111010000000
 ```

### IPV6 Unique local address
```
$ iptools ip6 describe -bits 64 -random -type unique-local
          Category                                  Value
---------------------------- ---------------------------------------------------
 IP Type                      Unique local
 Type Prefix                  fd00::/8
 IP                           fd00:1fc6:4f1e:9209:67dd:7bff:fedf:94c2
 Prefix                       fd00:1fc6:4f1e:9209::/64
 Global ID                    00:1fc6:4f1e
 Interface ID                 67dd:7bff:fedf:94c2
 Subnet ID                    9209
 Link                         http://[fd00:1fc6:4f1e:9209:67dd:7bff:fedf:94c2]/
 Subnet first address         fd00:1fc6:4f1e:9209:0000:0000:0000:0000
 Subnet last address          fd00:1fc6:4f1e:9209:ffff:ffff:ffff:ffff
 first address field binary   1111110100000000
```

### IPV6 multicast
```
$ iptools ip6 describe -bits 64 -random -type multicast
          Category                           Value
---------------------------- --------------------------------------
 IP Type                      Multicast
 Type Prefix                  ff00::/8
 IP                           ff1f:0:5419:9893:cba3:6592:2aa9:dc95
 Network Prefix               5419:9893:cba3:6592
 Group ID                     00:2aa9:dc95
 first address field binary   1111111100011111
```

### IPV6 Interface local multicast
```
$ iptools ip6 describe -bits 64 -random -type interface-local-multicast
          Category                           Value
---------------------------- --------------------------------------
 IP Type                      Interface local multicast
 Type Prefix                  ff00::/8
 IP                           ff31:0:3907:4492:6c5a:cc21:5397:1259
 Network Prefix               3907:4492:6c5a:cc21
 Group ID                     00:5397:1259
 first address field binary   1111111100110001
```

### IPV6 Link local multicast
```
$ iptools ip6 describe -bits 64 -random -type link-local-multicast
          Category                           Value
---------------------------- -------------------------------------
 IP Type                      Link local muticast
 Type Prefix                  ff00::/8
 IP                           ff02:0:3c06:5f2b:1959:e37d:5a4b:8f6
 Network Prefix               3c06:5f2b:1959:e37d
 Group ID                     00:5a4b:08f6
 first address field binary   1111111100000010
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
Go                              12            387            234           2033
Markdown                         1             52              0            315
-------------------------------------------------------------------------------
TOTAL                           13            439            234           2348
-------------------------------------------------------------------------------
```