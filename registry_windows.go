//go:build windows

package main

import "golang.org/x/sys/windows/registry"

func steamPathFromRegistry() string {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Valve\Steam`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return ""
	}
	defer key.Close()

	value, _, err := key.GetStringValue("SteamPath")
	if err != nil {
		return ""
	}

	return value
}
