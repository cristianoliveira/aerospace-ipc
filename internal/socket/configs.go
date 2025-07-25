package socket

import (
	"fmt"
	"os"

	"github.com/cristianoliveira/aerospace-ipc/internal/constants"
)

// GetSocketPath returns the socket path
//
// It checks for environment variable AEROSPACESOCK or uses the default socket path.
//
//	Default: /tmp/bobko.aerospace-<username>.sock
//	See: https://github.com/nikitabobko/AeroSpace/blob/f12ee6c9d914f7b561ff7d5c64909882c67061cd/Sources/Cli/_main.swift#L47
//
// Returns the socket path or an error if the path does not exist
func GetSocketPath() (string, error) {
	socketPath := fmt.Sprintf("/tmp/bobko.%s-%s.sock", "aerospace", os.Getenv("USER"))
	socketPathEnv := os.Getenv(constants.EnvAeroSpaceSock)
	if socketPathEnv != "" {
		socketPath = socketPathEnv
	}

	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return "", fmt.Errorf("failed to access socket path %s\r reason: %w", socketPath, err)
	}

	return socketPath, nil
}
