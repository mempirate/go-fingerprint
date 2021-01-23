# Go-fingerprint

Fingerprints devices on the network using crafted ARP requests and MAC OUIs.

Getting the fingerprint file:
```bash
sudo update-ieee-data # update ieee data (contains OUI file)
get-oui && cat ieee-oui.txt | grep -v '#' | sort > mac-fab.txt && rm ieee-oui.txt
```

Sources:
* [IEEE OUI list](http://standards-oui.ieee.org/oui/oui.txt)
* [arp-fingerprint](https://linux.die.net/man/1/arp-fingerprint)
    * Basically fuzzes the target with different payloads to see what responses are generated.

## Todo
* [Makefile](https://kodfabrik.com/journal/a-good-makefile-for-go)