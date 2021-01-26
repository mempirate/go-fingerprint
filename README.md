# Go-fingerprint

Fingerprints devices on the network using crafted ARP requests and MAC OUIs.

## Usage
**Linux**
```
$ go-fingerprint -i eth0 
2021/01/24 20:08:14
[*] Scanning on eth0: 172.18.8.216 [172.18.0.0/20]
IPv4                 MAC                  Hardware                      
===================================================================
172.18.0.1           00:15:5d:8f:xx:4c    Microsoft Corporation
172.18.0.3           28:c6:8e:xx:26:7b    NETGEAR             
172.18.0.3           04:18:d6:xx:ad:f7    Ubiquiti Networks Inc.
```
**Windows**
```
PS C:\> go-fingerprint.exe -i wi-fi
2021/01/24 20:08:14
[*] Scanning on \Device\NPF_{F27xxxxx-BF44-4xxx-8F50-2815xxxxxx85}: 192.168.1.103 [192.168.200.0/0]
IPv4                 MAC                  Hardware                      
===================================================================
192.168.1.8          04:18:d6:xx:ad:f7    Ubiquiti Networks Inc.
192.168.1.3          28:c6:8e:xx:26:7b    NETGEAR             
192.168.1.6          3c:2a:f4:xx:71:7f    Brother Industries, LTD.
192.168.1.11         40:cb:c0:xx:36:93    Apple, Inc.         
192.168.1.5          78:8a:20:xx:27:77    Ubiquiti Networks Inc.
192.168.1.102        28:3a:4d:xx:02:bb    Cloud Network Technology (Samoa) Limited
192.168.1.123        90:b0:ed:xx:46:cd    Apple, Inc.         
```


## Installation
Make sure your GOPATH environment variable is set and GOPATH/bin is added to your path.

**Linux & MacOS**

Dependencies on Linux:
* libpcap: `sudo apt-get install libpcap-dev`

Go get:
```
$ go get github.com/jonasbostoen/go-fingerprint
$ cd $GOPATH/src/github.com/jonasbostoen/go-fingerprint
$ make all
```

From source:
```
$ git clone https://github.com/jonasbostoen/go-fingerprint
$ cd go-fingerprint
$ make all
```

**Windows**

Go get:
```
PS C:\> go get github.com/jonasbostoen/go-fingerprint
```

## Todo
* Hardcode common OUIs
* Custom fingerprinting 
* Performance
    * https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
    * https://blog.golang.org/pprof

## Sources:
* [IEEE OUI list](http://standards-oui.ieee.org/oui/oui.txt)
* [arp-fingerprint](https://linux.die.net/man/1/arp-fingerprint)
    * Basically fuzzes the target with different payloads to see what responses are generated.

