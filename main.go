package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
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
}

var (
	iface = flag.String("i", "eth0", "Interface to listen on")
)

func main() {
	flag.Parse()
	scanner, err := GetInterface()
	if err != nil {
		log.Fatal(err)
	}
	err = ArpScan(scanner)
	if err != nil {
		log.Fatal(err)
	}
}

// GetInterface fills
func GetInterface() (*Interface, error) {
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
				// TODO: error checking:
				// * No loopback
				// * No localhost
				scanner = &Interface{
					iface:   i,
					ip:      ip4,
					netmask: ipnet.Mask,
				}
			}
		}
	}
	fmt.Printf("IP Address: %s %T\n", scanner.ip, scanner.ip)
	fmt.Printf("Network: %s, netmask: %s\n", scanner.ip.Mask(scanner.netmask), scanner.netmask)
	return scanner, nil
}

// ArpScan scans the network using the interface provided
func ArpScan(scanner *Interface) error {
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

	// Start sending ARP requests

	for {
		for _, ip := range GetIPAddresses(scanner.ip, scanner.netmask) {
			arp.DstProtAddress = []byte(ip)
			gopacket.SerializeLayers(buf, opts, &eth, &arp)
			if err := handle.WritePacketData(buf.Bytes()); err != nil {
				return err
			}
		}
		time.Sleep(time.Second * 10)
	}

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
			// Note:  we might get some packets here that aren't responses to ones we've sent,
			// if for example someone else sends US an ARP request.  Doesn't much matter, though...
			// all information is good information :)
			log.Printf("IP %v is at %v", net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
		}
	}
}

// GetIPAddresses returns all IP addresses on a subnet
func GetIPAddresses(ip net.IP, mask net.IPMask) (out []net.IP) {
	bip := binary.BigEndian.Uint32([]byte(ip))
	bmask := binary.BigEndian.Uint32([]byte(mask))
	bnet := bip & bmask
	bbroadcast := bnet | ^bmask

	for bnet++; bnet < bbroadcast; bnet++ {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], bnet)
		out = append(out, net.IP(buf[:]))
	}
	return
}
