# Go-fingerprint

Fingerprints devices on the network using crafted ARP requests and MAC OUIs.

## Installation
Install with Go on Linux, Windows and Mac:
`go get github.com/jonasbostoen/go-fingerprint`

Make sure your GOPATH environment variable is set and GOPATH/bin is added to your path.

**Linux**
Dependencies:
* libpcap: `sudo apt-get install libpcap-dev`

To avoid having to use `sudo`, give the binary the following capability:

`sudo setcap CAP_NET_RAW+ep go-fingerprint`

Sources:
* [IEEE OUI list](http://standards-oui.ieee.org/oui/oui.txt)
* [arp-fingerprint](https://linux.die.net/man/1/arp-fingerprint)
    * Basically fuzzes the target with different payloads to see what responses are generated.

## Todo
* Hardcode common OUIs
