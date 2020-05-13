package main

import (
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
)

type PppPAPCode uint8

const (
	AuthenticateRequest PppPAPCode = 0x01
	//AuthenticateAck     PppPAPCode = 0x02
	//AuthenticateNack    PppPAPCode = 0x03
)

type PppPAP struct {
	code       PppPAPCode
	identifier uint8
	length     uint16
	peerID     []byte
	password   []byte
}

func DecodePppPAP(ppp *layers.PPP) *PppPAP {
	peerIdLength := ppp.Payload[4]
	passwordLength := ppp.Payload[5+peerIdLength]
	return &PppPAP{
		code:       PppPAPCode(ppp.Payload[0]),
		identifier: ppp.Payload[1],
		length:     binary.BigEndian.Uint16(ppp.Payload[2:4]),
		peerID:     ppp.Payload[5 : 5+peerIdLength],
		password:   ppp.Payload[6+peerIdLength : 6+peerIdLength+passwordLength],
	}
}

func PrintPapInfoIfPossible(packet gopacket.Packet) {
	pppLayer := packet.Layer(layers.LayerTypePPP)

	if pppLayer == nil {
		return
	}

	ppp, _ := pppLayer.(*layers.PPP)

	if ppp.PPPType != 0xc023 {
		return
	}

	pap := DecodePppPAP(ppp)

	if pap.code != AuthenticateRequest {
		return
	}

	log.Printf(`

=================================================
===  PAP Username : %s
===  PAP Password : %s
=================================================

`, string(pap.peerID), string(pap.password))
}
