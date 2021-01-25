package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/bits"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Interface represents an interface
type Interface struct {
	iface   *net.Interface
	ip      net.IP
	netmask net.IPMask
	prefix  uint8
}

var (
	iface = flag.String("i", "eth0", "Interface to listen on")
)

func main() {
	flag.Parse()
	scanner, err := getInterface()
	if err != nil {
		log.Fatal(err)
	}
	err = arpScan(scanner)
	if err != nil {
		log.Fatal(err)
	}
}

/// Gets interface based on flag (or default ethernet)
func getInterface() (*Interface, error) {
	i, err := net.InterfaceByName(*iface)

	if err != nil {
		return nil, errors.New("Can't get interface " + *iface)
	}

	var scanner *Interface

	// TODO: account for multiple ipv4s on 1 interface
	addrs, err := i.Addrs()

	if err != nil {
		return nil, errors.New("Problem getting interface addresses")
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				scanner = &Interface{
					iface:   i,
					ip:      ip4,
					netmask: ipnet.Mask,
					prefix:  uint8(bits.OnesCount32(binary.BigEndian.Uint32(ipnet.Mask))),
				}
			}
		}
	}
	return scanner, nil
}

// arpScan scans the network using the interface provided
func arpScan(scanner *Interface) error {
	handle, err := pcap.OpenLive(scanner.iface.Name, 1024, false, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Start reading ARP packets in a goroutine
	stop := make(chan struct{})
	go readARP(handle, scanner.iface, stop)
	defer close(stop)

	// Set up the layers
	eth := layers.Ethernet{
		SrcMAC:       scanner.iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(scanner.iface.HardwareAddr),
		SourceProtAddress: []byte(scanner.ip),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	log.Printf("\n[*] Scanning on %s: %s [%s/%d]\n", scanner.iface.Name, scanner.ip, scanner.ip.Mask(scanner.netmask), scanner.prefix)
	fmt.Printf("%-20s %-20s %-30s\n", "IPv4", "MAC", "Hardware")
	fmt.Println("===================================================================")

	// Start sending ARP requests
	for _, ip := range getIPAddresses(&scanner.ip, &scanner.netmask) {
		arp.DstProtAddress = []byte(ip)
		gopacket.SerializeLayers(buf, opts, &eth, &arp)
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return err
		}
	}

	// Wait for ARP responses (tune this to network size)
	time.Sleep(time.Second * 3)

	return nil
}

func readARP(handle *pcap.Handle, iface *net.Interface, stop chan struct{}) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

	for {
		var packet gopacket.Packet
		select {
		case <-stop:
			return
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
				// This is a packet I sent.
				continue
			}

			go examineMAC(arp.SourceProtAddress, arp.SourceHwAddress)
		}
	}
}

func examineMAC(ip, mac []byte) {
	var fabmatch string

	oui := mac[:3]

	f, err := os.Open("mac-fab.txt")
	if err != nil {
		fmt.Println("Can't find mac-fab.txt, continuing without fingerprinting")
	}

	defer f.Close()

	input := bufio.NewScanner(f)
	for input.Scan() {
		line := strings.Fields(input.Text())
		macstr := line[0]
		fab := strings.Join(line[1:], " ")
		macbytes, err := hex.DecodeString(macstr)
		if err != nil {
			fmt.Println(err)
		}

		if bytes.Compare(oui, macbytes) == 0 {
			fabmatch = fab
		}
	}
	fmt.Printf("%-20v %-20v %-20s\n", net.IP(ip), net.HardwareAddr(mac), fabmatch)

}

// getIPAddresses returns all IP addresses on a subnet
func getIPAddresses(ip *net.IP, mask *net.IPMask) (out []net.IP) {
	bip := binary.BigEndian.Uint32([]byte(*ip))
	bmask := binary.BigEndian.Uint32([]byte(*mask))
	bnet := bip & bmask
	bbroadcast := bnet | ^bmask

	for bnet++; bnet < bbroadcast; bnet++ {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], bnet)
		out = append(out, net.IP(buf[:]))
	}
	return
}
