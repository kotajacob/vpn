package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
)

const cacheDirName = "vpn"

type config struct {
	name  string
	usage int
}

type configs []config

func (c configs) Len() int      { return len(c) }
func (c configs) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c configs) Less(i, j int) bool {
	if c[i].usage == c[j].usage {
		// Sort alphabetically as a backup.
		return c[i].name < c[j].name
	}
	return c[i].usage > c[j].usage
}

// recent returns a list of interfaces sorted by how recently they've been used.
func recent() (configs, error) {
	var configs configs
	usage := getUsage()

	entries, err := os.ReadDir(wireguardConfigs)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot read wireguard config dir %v: %v\n",
			wireguardConfigs,
			err,
		)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		configs = append(configs, config{
			name:  entry.Name(),
			usage: usage[entry.Name()],
		})
	}
	sort.Sort(configs)
	return configs, nil
}

// getUsage builds a map of wireguard configs and their usage values.
func getUsage() map[string]int {
	usage := make(map[string]int)
	dir := filepath.Join(xdg.CacheHome, cacheDirName)
	if !exists(filepath.Join(xdg.CacheHome, cacheDirName)) {
		return usage
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return usage
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		b, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		used, err := strconv.Atoi(strings.TrimSpace(string(b)))
		if err != nil {
			continue
		}
		usage[entry.Name()] = used
	}
	return usage
}

// exists returns true if the specified path exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
