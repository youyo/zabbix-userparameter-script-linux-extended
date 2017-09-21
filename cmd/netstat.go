package cmd

import (
	"errors"
	"fmt"

	"github.com/drael/GOnetstat"
	"github.com/spf13/cobra"
)

type (
	Netstat struct {
		State    string
		Protocol string
	}
)

var netstat Netstat

// netstatCmd represents the netstat command
var netstatCmd = &cobra.Command{
	Use: "netstat",
	//Short: "A brief description of your command",
	//Long: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		value, err := netstatCount()
		if err != nil {
			fmt.Println(err)
			fmt.Println("ZBX_NOTSUPPORTED")
		} else {
			fmt.Println(value)
		}
	},
}

func netstatCount() (value uint64, err error) {
	state := "LISTEN"
	protocol := "tcp"
	if netstat.State != "" {
		state = netstat.State
	}
	if netstat.Protocol != "" {
		protocol = netstat.Protocol
	}

	var d []GOnetstat.Process
	switch protocol {
	case "tcp":
		d = GOnetstat.Tcp()
	case "udp":
		d = GOnetstat.Udp()
	case "tcp6":
		d = GOnetstat.Tcp6()
	case "udp6":
		d = GOnetstat.Udp6()
	default:
		err = errors.New("Not match protocol")
		return
	}
	switch state {
	case
		"LISTEN",
		"ESTABLISHED",
		"TIME_WAIT",
		"LISTENING",
		"SYN_SENT",
		"SYN_RECEIVED",
		"FIN_WAIT_1",
		"FIN_WAIT_2",
		"CLOSE_WAIT",
		"CLOSING",
		"LAST_ACK",
		"CLOSED":
		value = countUp(d, state)
	default:
		err = errors.New("Not match State")
		return
	}
	return
}

func countUp(data []GOnetstat.Process, state string) (value uint64) {
	for _, v := range data {
		if v.State == state {
			value++
		}
	}
	return
}

func init() {
	RootCmd.AddCommand(netstatCmd)
	netstatCmd.Flags().StringVarP(&netstat.State, "state", "s", "", "State")
	netstatCmd.Flags().StringVarP(&netstat.Protocol, "protocol", "p", "", "Protocol")
}
