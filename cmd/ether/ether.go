package ether

import (
	"github.com/QQGoblin/ssconnect/pkg/ether"
	"github.com/spf13/cobra"
)

var (
	iface        string
	ethernetType uint16
	dst          string
	src          string
)

func init() {

	Command.PersistentFlags().StringVarP(&iface, "iface", "i", "eth0", "bind socket to interface")
	Command.PersistentFlags().Uint16Var(&ethernetType, "etherType", 0xCCCC, "custom ether type define")

	SendCMD.PersistentFlags().StringVarP(&dst, "dst", "d", "ff:ff:ff:ff:ff:ff", "destination mac")
	SendCMD.PersistentFlags().StringVarP(&src, "src", "s", "aa:bb:cc:11:11:11", "source mac")

	Command.AddCommand(ReceiveCMD)
	Command.AddCommand(SendCMD)
}

var Command = &cobra.Command{
	Use:   "ether",
	Short: "send and receiver message by ethernet frame",
}

var ReceiveCMD = &cobra.Command{
	Use:   "receive",
	Short: "receive message",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ether.ReceiveEtherPkt(iface, ethernetType)
	},
}

var SendCMD = &cobra.Command{
	Use:   "send",
	Short: "send message",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ether.SendEtherPkt(src, dst, iface, ethernetType, []byte(args[0]))
	},
}
