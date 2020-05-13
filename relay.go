package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"net"
)

type CaptureInfo struct {
	fromIface      *net.Interface
	fromHandle     *pcap.Handle
	toIface        *net.Interface
	toHandle       *pcap.Handle
	macConnections *cmap.ConcurrentMap
}

func Bridge(from *net.Interface, to *net.Interface, macConnections *cmap.ConcurrentMap) {
	fromHandle, err := pcap.OpenLive(from.Name, 65536, true, pcap.BlockForever)

	if err != nil {
		log.Fatalf("Error while connecting to the interface %s: %s\n", from.Name, err.Error())
	}

	toHandle, err := pcap.OpenLive(to.Name, 65536, true, pcap.BlockForever)

	if err != nil {
		log.Fatalf("Error while connecting to the interface %s: %s\n", to.Name, err.Error())
	}

	captureInfo := CaptureInfo{
		fromIface:      from,
		fromHandle:     fromHandle,
		toIface:        to,
		toHandle:       toHandle,
		macConnections: macConnections,
	}

	packetSource := gopacket.NewPacketSource(fromHandle, fromHandle.LinkType())

	handle(&captureInfo, packetSource)
}

func doNotSendPacketBack(packet gopacket.Packet, info *CaptureInfo) bool {
	if layer := packet.Layer(layers.LayerTypeEthernet); layer != nil {
		layer, _ := layer.(*layers.Ethernet)

		info.macConnections.SetIfAbsent(layer.SrcMAC.String(), info.fromIface.Name)

		ifaceName, ok := info.macConnections.Get(layer.SrcMAC.String())

		/*
			to avoid circular packet sending between interfaces,
			do not send packets coming from devices this interface is connected with
		*/
		return ok && ifaceName != info.fromIface.Name
	}

	return false
}

func handle(info *CaptureInfo, packetSource *gopacket.PacketSource) {
	for packet := range packetSource.Packets() {

		if doNotSendPacketBack(packet, info) {
			continue
		}

		go PrintPapInfoIfPossible(packet)

		err := info.toHandle.WritePacketData(packet.Data())

		if err != nil {
			log.Printf("Error while sending packet, %s\n", err.Error())
		}
	}
}
