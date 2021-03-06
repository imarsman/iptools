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

One thing I'd like to do is try IP6 subnetting.

### Subnet divisions split into non-default sized networks

```
$ iptools subnet divide -ip 99.236.32.0 -bits 16 -secondary-bits 18 -pretty


      Category            Value
-------------------- ---------------
 Subnet               99.236.0.0/16
 Secondary Subnet     99.236.0.0/18
 Networks             1
 Secondary Networks   4
 Effective Networks   4
 Network Hosts        65536
 Sub Network Hosts    16384

     Subnet
-----------------
 99.236.0.0/18
 99.236.64.0/18
 99.236.128.0/18
 99.236.192.0/18
 ```

### Subnet divisions split default

```
$ iptools subnet divide -ip 99.236.32.0 -bits 16 -pretty


   Category          Value
--------------- ---------------
 Subnet          99.236.0.0/16
 Networks        1
 Network Hosts   65536

    Subnet
---------------
 99.236.0.0/16
 ```

### Subnet ranges in default size

```
$ iptools subnet ranges -ip 99.236.32.0 -bits 16 -pretty


   Category          Value
--------------- ---------------
 Subnet          99.236.0.0/16
 Networks        1
 Network Hosts   65536

   Start           End
------------ ----------------
 99.236.0.0   99.236.255.255
 ```
 
### Subnet ranges split into non-default size

```
$ iptools subnet ranges -ip 99.236.32.0 -bits 16 -secondary-bits 18 -pretty


      Category            Value
-------------------- ---------------
 Subnet               99.236.0.0/16
 Secondary Subnet     99.236.0.0/18
 Networks             1
 Secondary Networks   4
 Effective Networks   4
 Network Hosts        65536
 Sub Network Hosts    16384

    Start            End
-------------- ----------------
 99.236.0.0     99.236.63.255
 99.236.64.0    99.236.127.255
 99.236.128.0   99.236.191.255
 99.236.192.0   99.236.255.255
```


### Subnet details

```
$ iptools subnet describe -ip 10.32.0.0 -bits 24
+--------------------+-------------------------------------+
|      Category      |                Value                |
+--------------------+-------------------------------------+
| Subnet             | 10.32.0.0/24                        |
| Network Address    | 10.32.0.0                           |
| IP Address         | 10.32.0.255                         |
| Broadcast Address  | 10.32.0.255                         |
| Networks           | 1                                   |
| Network Hosts      | 256                                 |
| Total Hosts        | 256                                 |
| IP Class           | A                                   |
| IP Type            | Private                             |
| Binary Subnet Mask | 00001010.00100000.00000000.00000000 |
| Binary ID          | 00001010001000000000000000000000    |
| Hex ID             | a2000ff                             |
| in-addr.arpa       | 0.0.32.10.in-addr.arpa              |
| Wildcard Mask      | 240.218.250.250                     |
+--------------------+-------------------------------------+
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
  subnet                 Get networks for subnet
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

### Lines of cde

```
$ gocloc pkg/subnet pkg/util cmd README.md
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                               9            207            145           1189
Markdown                         1             35              0            128
-------------------------------------------------------------------------------
TOTAL                           10            242            145           1317
-------------------------------------------------------------------------------
```