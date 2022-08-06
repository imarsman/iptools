# iptools

This project was done as a learning experience for things like IP subnetting. Mostly I was interested in learning more
about ipv4 subnets and splitting up subnets into equal ranges and defining things like the network ID and broadcast
address for subnets. This is a way for me to better establish my knowledge of this subject.

This utility can have completion enabled by typing `COMP_INSTALL=1 iptools`.

This utility currently does three things. It splits a subnet into networks and into networks by a differing subnet size,
it splits a subnet into a set of ranges for its networks, and it gives summary information for a subnet.

The IP package used, `net/netip`, is in Go 1.18. It is slightly different in API compared to the
[netaddr](https://github.com/inetaf/netaddr) package that came first. Sadly, the `netip` package loses the IPRange
struct, so I have added needed functionality here in the subnet package.

One thing I'd like to do is try IP6 subnetting. I have to learn more about IPV6 and its subnetting first.

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
 Network Hosts       65536

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
 Secondary Subnet                     99.236.0.0/18
 Secondary Subnet IP                  99.236.0.0
 Secondary Subnet Broadcast Address   99.236.255.255
 Secondary Subnet Mask                255.255.192.0
 Networks                             1
 Secondary Networks                   4
 Effective Networks                   4
 Network Hosts                        65536
 Sub Network Hosts                    16384

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
      Category                       Value
-------------------- -------------------------------------
 Subnet               10.32.0.0/23
 Subnet IP            10.32.0.0
 Broadcast Address    10.32.1.255
 Subnet Mask          255.255.254.0
 Wildcard Mask        0.0.1.255
 IP Class             A
 IP Type              Private
 Binary Subnet Mask   00001010.00100000.00000000.00000000
 Binary ID            00001010001000000000000000000000
 Hex ID               0xA2001FF
 in-addr.arpa         0.0.32.10.in-addr.arpa
 Networks             128
 Network Hosts        512
```

#### Subnet with secondary subnet details

```
$ iptools subnetip4 describe -ip 10.32.0.0 -bits 16 -secondary-bits 18
              Category                               Value
------------------------------------ -------------------------------------
 Subnet                               10.32.0.0/16
 Subnet IP                            10.32.0.0
 Broadcast Address                    10.32.255.255
 Subnet Mask                          255.255.0.0
 Wildcard Mask                        0.0.255.255
 IP Class                             A
 IP Type                              Private
 Binary Subnet Mask                   00001010.00100000.00000000.00000000
 Binary ID                            00001010001000000000000000000000
 Hex ID                               0xA20FFFF
 in-addr.arpa                         0.0.32.10.in-addr.arpa
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
```

### Top level help

```
$ iptools -h
iptools
-------
Commit:  fe61dd8
Date:    2022-07-12T02:59:52Z
Tag:     v0.1.11
OS:      darwin
ARCH:    arm64

Usage: iptools <command> [<args>]

Options:
  --help, -h             display this help and exit
  --version              display version and exit

Commands:
  subnetip4              Get networks for subnet
```

### Help for subnet options

```
$ iptools subnet -h
iptools
-------
Commit:  fe61dd8
Date:    2022-07-12T02:59:52Z
Tag:     v0.1.11
OS:      darwin
ARCH:    arm64

Usage: iptools subnet <command> [<args>]
  --help, -h             display this help and exit
  --version              display version and exit

Commands:
  ranges                 divide a subnet into ranges
  divide                 divide a subnet into smaller subnets
  describe               describe a subnet
```

### Lines of code

```
$ gocloc pkg cmd README.md
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                              11            301            229           1533
Markdown                         1             37              0            187
-------------------------------------------------------------------------------
TOTAL                           12            338            229           1720
-------------------------------------------------------------------------------
```