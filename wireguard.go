package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
)

// active checks for an active wireguard vpn connection and returns the config
// if found.
func active(configs configs) (bool, config, error) {
	cmd := exec.Command("ip", "a")
	b, err := cmd.Output()
	if err != nil {
		return false, config{}, err
	}

	// Find lines like the following:
	// 6: nz-akl-wg-301: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1420 qdisc noqueue state UNKNOWN group default qlen 1000
	for _, line := range strings.Split(string(b), "\n") {
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}

		// Check if the line begins with a number and colon.
		if !strings.HasSuffix(parts[0], ":") {
			continue
		}
		if _, err := strconv.Atoi(strings.TrimSuffix(parts[0], ":")); err != nil {
			continue
		}

		// Check if device matches a wireguard config.
		if !strings.HasSuffix(parts[1], ":") {
			continue
		}
		device := strings.TrimSuffix(parts[1], ":")
		for _, c := range configs {
			if device+".conf" == c.name {
				return true, c, nil
			}
		}
	}
	return false, config{}, nil
}

// enable a wireguard config and increases usage value.
func enable(c config) error {
	// TODO: Support sudo instead of only doas.
	cmd := exec.Command(
		"doas",
		"wg-quick",
		"up",
		strings.TrimSuffix(c.name, ".conf"),
	)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to enable %v: %v", c.name, err)
	}

	// Update usage stats.
	err = os.MkdirAll(xdg.CacheHome, 0o744)
	if err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}
	err = os.MkdirAll(filepath.Join(xdg.CacheHome, cacheDirName), 0o700)
	if err != nil {
		return fmt.Errorf("failed to create usage cache directory: %v", err)
	}
	err = os.WriteFile(
		filepath.Join(xdg.CacheHome, cacheDirName, c.name),
		[]byte(strconv.Itoa(c.usage+1)),
		0o600,
	)
	if err != nil {
		return fmt.Errorf("failed to update usage: %v", err)
	}
	return nil
}

// disable a wireguard config.
func disable(c config) error {
	// TODO: Support sudo instead of only doas.
	cmd := exec.Command(
		"doas",
		"wg-quick",
		"down",
		strings.TrimSuffix(c.name, ".conf"),
	)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to enable %v: %v", c.name, err)
	}
	return nil
}
