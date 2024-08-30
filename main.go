package main

import (
	"fmt"
	"github.com/QQGoblin/ssconnect/cmd/ether"
	"github.com/QQGoblin/ssconnect/cmd/ipproto"
	"github.com/QQGoblin/ssconnect/cmd/udp"
	"github.com/spf13/cobra"
	"os"
)

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "niptalk",
		Short: "this is healp tools for connect on local network",
	}

	cmd.AddCommand(udp.Command)
	cmd.AddCommand(ipproto.Command)
	cmd.AddCommand(ether.Command)
	return cmd

}

func main() {

	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Printf("exit with error, %+v\n", err)
		os.Exit(1)
	}
}
