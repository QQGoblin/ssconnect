package ether

import (
	"fmt"
	"github.com/QQGoblin/ssconnect/pkg/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"
	"net"
	"os"
	"syscall"
)

// bindToInterface 将 socket 绑定到指定设备
func bindToInterface(fd int, ifaceName string, ethernetType uint16) error {
	ifIndex := 0
	// An empty string here means to listen to all interfaces
	if ifaceName != "" {
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			return fmt.Errorf("InterfaceByName: %v", err)
		}
		ifIndex = iface.Index
	}
	s := &syscall.SockaddrLinklayer{
		Protocol: utils.Htons(ethernetType),
		Ifindex:  ifIndex,
	}
	return syscall.Bind(fd, s)
}

func ReceiveEtherPkt(ifaceName string, ethernetType uint16) error {

	// 创建套接字
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(utils.Htons(ethernetType)))
	if err != nil {
		return errors.Wrapf(err, "syscall socket failed")
	}
	fmt.Printf("Obtained fd %d and listen for Ethernet(0x%x)\n", fd, int(ethernetType))
	defer syscall.Close(fd)

	// 绑定网卡
	if err = bindToInterface(fd, ifaceName, ethernetType); err != nil {
		return errors.Wrapf(err, "bind socket to %s failed", ifaceName)
	}

	// 初始化解码工具
	var ethLayer layers.Ethernet
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &ethLayer)
	parser.IgnoreUnsupported = true
	decodedLayers := make([]gopacket.LayerType, 0)

	// 接收数据
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	defer f.Close()

	for {
		buf := make([]byte, 2048)

		if _, err = f.Read(buf); err != nil {
			return errors.Wrapf(err, "problems @ location 2")
		}

		// 解析从原始报文中解析报文头
		if err = parser.DecodeLayers(buf, &decodedLayers); err != nil {
			return errors.Wrapf(err, "problems parase ethLayer header")
		}

		fmt.Printf("Receive Ethernet(0x%x) %s -> %s: %s\n", int(ethLayer.EthernetType), ethLayer.SrcMAC.String(), ethLayer.DstMAC.String(), string(ethLayer.Payload))
	}
}

func SendEtherPkt(src, dst, ifaceName string, ethernetType uint16, payload []byte) error {

	// 构建 ethernet packet
	h, err := utils.BuildEtherHeader(src, dst, ethernetType)
	if err != nil {
		return errors.Wrapf(err, "Failed to build ethernet header")
	}

	buf, err := utils.MarshalPkt(payload, h)
	if err != nil {
		return errors.Wrapf(err, "Failed to Marshal Ethernet Packet")
	}

	// 创建套接字
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(utils.Htons(ethernetType)))
	if err != nil {
		return errors.Wrapf(err, "Failed to create raw socket")
	}
	defer syscall.Close(fd)

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return fmt.Errorf("InterfaceByName: %v", err)
	}

	// 发送数据到指定网卡
	addr := &syscall.SockaddrLinklayer{
		Protocol: utils.Htons(ethernetType),
		Ifindex:  iface.Index,
	}

	// 参考 https://man7.org/linux/man-pages/man3/sendto.3p.html
	return syscall.Sendto(fd, buf, 0, addr)
}
