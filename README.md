# iptools

This is a learning experience for things like IP subnetting. Mostly I am interested in splitting up subnets into equal
ranges and defining things like the network ID and broadcast address for subnets. This is a way for me to better
establish my knowledge of this subject.

This utility can have completion enabled by typing `COMP_INSTALL=1 iptools`.

This utility currently does three things. It splits a subnet into networks and into networks by a differing subnet size,
it splits a subnet into a set of ranges for its networks, and it gives summary information for a subnet.

Subnet networks

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

Subnet ranges

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


Subnet details
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