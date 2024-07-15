package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const wireguardConfigs = "/etc/wireguard"

func main() {
	flag.Parse()
	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		usage()
	}

	configs, err := recent()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read configs:", err)
		os.Exit(1)
	}

	isActive, activeConfig, err := active(configs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed while detecting active vpn connection:", err)
		os.Exit(1)
	}

	if len(flag.Args()) == 1 {
		switch flag.Arg(0) {
		case "on", "up":
			fmt.Println("enabling", configs[0].name)
			if err := enable(configs[0]); err != nil {
				fmt.Fprintln(os.Stderr, "failed enabling vpn:", err)
				os.Exit(1)
			}
		case "off", "down":
			if !isActive {
				fmt.Println("vpn is not active")
				os.Exit(0)
			}
			fmt.Println("disabling", activeConfig.name)
			if err := disable(activeConfig); err != nil {
				fmt.Fprintln(os.Stderr, "failed enabling vpn:", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintln(os.Stderr, "unknown argument", flag.Arg(0))
			os.Exit(1)
		}
	} else {
		selected := -1
		if isActive {
			for i, c := range configs {
				if c.name == activeConfig.name {
					selected = i
				}
			}
		}
		p := tea.NewProgram(model{
			choices:  configs,
			selected: selected,
		}, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "failed running TUI:", err)
			os.Exit(1)
		}
	}
}

// usage prints usage information and exits with an error.
func usage() {
	fmt.Fprintln(os.Stderr, "usage: vpn [ on | up | off | down ]")
	os.Exit(1)
}
