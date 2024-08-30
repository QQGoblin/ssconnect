package utils

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"
	"net"
)

func Htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func MarshalPkt(payload []byte, header ...gopacket.SerializableLayer) ([]byte, error) {
	buf := gopacket.NewSerializeBuffer()
	var layersToSerialize []gopacket.SerializableLayer

	layersToSerialize = append(layersToSerialize, header...)
	layersToSerialize = append(layersToSerialize, gopacket.Payload(payload))

	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}, layersToSerialize...); err != nil {
		return nil, errors.Wrapf(err, "Error serializing packet")
	}

	return buf.Bytes(), nil
}

func BuildEtherHeader(src, dst string, ethernetType uint16) (*layers.Ethernet, error) {
	srcMAC, err := net.ParseMAC(src)
	if err != nil {
		return nil, errors.Wrapf(err, "error src %s", src)
	}

	dstMAC, err := net.ParseMAC(dst)
	if err != nil {
		return nil, errors.Wrapf(err, "error dst %s", src)
	}

	return &layers.Ethernet{
		BaseLayer:    layers.BaseLayer{},
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetType(Htons(ethernetType)),
	}, nil
}

func BuildIPv4Header(src, dst string, ipproto uint8) (*layers.IPv4, error) {
	return &layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    net.ParseIP(src),
		DstIP:    net.ParseIP(dst),
		Protocol: layers.IPProtocol(ipproto),
	}, nil
}
