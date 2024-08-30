package ipproto

import (
	"fmt"
	"github.com/QQGoblin/ssconnect/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"syscall"
)

func SendIPPkt(srcIP, dstIP string, protocol uint8, payload []byte) error {

	// 生成 IP 数据包
	h, err := utils.BuildIPv4Header(srcIP, dstIP, protocol)
	if err != nil {
		return errors.Wrapf(err, "Failed to build IPv4 Header")
	}

	buf, err := utils.MarshalPkt(payload, h)
	if err != nil {
		return errors.Wrapf(err, "Failed to Marshal IPv4 Packet")
	}

	// 打开套接字
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return errors.Wrapf(err, "Failed to create raw socket")
	}
	defer syscall.Close(fd)

	// 可以通过设置`IP_HDRINCL`参数自定义网络包的`IP Header`信息（PS：未指定该参数时，系统会自动填充`IP`头）
	if err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		return errors.Wrapf(err, "Failed to set IP_HDRINCL")
	}

	dstAddr := net.ParseIP(dstIP).To4()
	addr := &syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{dstAddr[0], dstAddr[1], dstAddr[2], dstAddr[3]},
	}

	// 参考 https://man7.org/linux/man-pages/man3/sendto.3p.html
	return syscall.Sendto(fd, buf, 0, addr)
}

func ReceiveIPPkt(protocol uint8) error {

	// 打开套接字
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, int(protocol))
	if err != nil {
		return errors.Wrapf(err, "syscall socket failed")
	}
	fmt.Printf("Obtained fd %d\n", fd)
	defer syscall.Close(fd)

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	defer f.Close()

	for {
		buf := make([]byte, 1500) // Base on mtu
		numRead, err := f.Read(buf)
		if err != nil {
			return errors.Wrapf(err, "problems @ location 2")
		}

		// 解析从原始报文中解析 IP 报文头
		ipHeader, err := ipv4.ParseHeader(buf[:numRead])
		if err != nil {
			return errors.Wrapf(err, "problems parase ipv4 header")
		}

		fmt.Printf("Receive %s -> %s: %s\n", ipHeader.Src.String(), ipHeader.Dst.String(), buf[ipv4.HeaderLen:])
	}
}
