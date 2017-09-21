package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youyo/zabbix-userparameter-script-linux-extended/libs/linux-extended"
)

const swapFile = "/proc/swaps"

type (
	Swap struct {
		Action string
		Device string
		Unit   string
	}
)

var swap Swap

var swapCmd = &cobra.Command{
	Use: "swap",
	//Short: "A brief description of your command",
	//Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		switch swap.Action {
		case "discovery":
			d, _ := swapDiscovery()
			fmt.Println(d.Json())
		case "size":
			value, err := swapSize()
			if err != nil {
				fmt.Println(err)
				fmt.Println("ZBX_NOTSUPPORTED")
			} else {
				fmt.Println(value)
			}
		default:
			fmt.Println("ZBX_NOTSUPPORTED")
		}
	},
}

func swapDiscovery() (d linux_extended.DiscoveryData, err error) {
	f, err := os.Open(swapFile)
	defer f.Close()
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		if !strings.HasPrefix(s.Text(), "Filename") {
			device := strings.Fields(s.Text())[0]
			d = append(d, linux_extended.DiscoveryItem{
				"DEVICE": device,
			})
		}
	}
	return
}

func swapSize() (value uint64, err error) {
	device := ""
	unit := "used"
	if swap.Device != "" {
		device = swap.Device
	}
	if swap.Unit != "" {
		unit = swap.Unit
	}
	var l []string
	f, err := os.Open(swapFile)
	defer f.Close()
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), device) {
			l = strings.Fields(s.Text())
		}
	}
	total, err := strconv.ParseUint(l[2], 10, 64)
	if err != nil {
		return
	}
	used, err := strconv.ParseUint(l[3], 10, 64)
	if err != nil {
		return
	}

	switch unit {
	case "used":
		value = used
	case "total":
		value = total
	case "free":
		value = total - used
	case "pfree":
		value = (total - used) * 100 / total
	}
	return
}

func init() {
	RootCmd.AddCommand(swapCmd)
	swapCmd.Flags().StringVarP(&swap.Action, "action", "a", "", "Action")
	swapCmd.Flags().StringVarP(&swap.Device, "device", "d", "", "Device")
	swapCmd.Flags().StringVarP(&swap.Unit, "unit", "u", "", "Unit")
}
