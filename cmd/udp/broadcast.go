package udp

import (
	"github.com/QQGoblin/ssconnect/pkg/udp"
	"github.com/spf13/cobra"
)

var (
	bPort int
)

func init() {

	Command.PersistentFlags().IntVarP(&bPort, "port", "p", 18789, "udp listen port")
	Command.AddCommand(BroadcastCMD)
	Command.AddCommand(ReceiveCMD)
}

var Command = &cobra.Command{
	Use:   "udp",
	Short: "send and receiver message by udp broadcast",
}

var BroadcastCMD = &cobra.Command{
	Use:     "broadcast",
	Aliases: []string{"bd"},
	Short:   "send message",
	RunE: func(cmd *cobra.Command, args []string) error {
		return udp.Broadcast(bPort, args[0])
	},
}

var ReceiveCMD = &cobra.Command{
	Use:   "receive",
	Short: "receive message",
	RunE: func(cmd *cobra.Command, args []string) error {
		return udp.Receive(bPort)
	},
}
