package client

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/internal/socket"
)

// Connector should return AeroSpaceConnectiono
// AeroSpaceConnector is an interface for connecting to the AeroSpace socket.
//
// It provides a method to establish a connection and return an AeroSpaceSocketConn.
//
//	See: AeroSpaceDefaultConnector for the default implementation.
//
// It allows one to set their custom connector if needed, for testing or other purposes.
type AeroSpaceConnector interface {
	// Connect to the AeroSpace Socket and return client
	Connect() (AeroSpaceSocketConn, error)
}

// AeroSpaceDefaultConnector is the default implementation of AeroSpaceConnector.
//
// In most cases, you will use this connector to connect to the AeroSpace socket.
type AeroSpaceDefaultConnector struct{}

func (c *AeroSpaceDefaultConnector) Connect() (AeroSpaceSocketConn, error) {
	socketPath, err := socket.GetSocketPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get socket path\n %w", err)
	}

	client, err := NewAeroSpaceSocketConnection(socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to creat socket connection\n%w", err)
	}

	return client, nil
}

// AeroSpaceCustomConnector is the default implementation of AeroSpaceConnector.
//
// In most cases, you will use this connector to connect to the AeroSpace socket.
type AeroSpaceCustomConnector struct {
	// SocketPath is the custom socket path for the AeroSpace connection.
	SocketPath string
	// ValidateVersion indicates whether to validate the version of the AeroSpace server.
	ValidateVersion bool
}

// Connect establishes a connection to the AeroSpace socket and validates the server version
// with the minimum required version.
//
// It returns an AeroSpaceSocketConn or an error if the connection fails.
func (c *AeroSpaceCustomConnector) Connect() (AeroSpaceSocketConn, error) {
	if c.SocketPath == "" {
		return nil, fmt.Errorf("socket path cannot be empty")
	}

	client, err := NewAeroSpaceSocketConnection(c.SocketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to creat socket connection\n%w", err)
	}

	response, err := client.SendCommand("config", []string{"--config-path"})
	if err != nil {
		return nil, fmt.Errorf("failed communicate with server\n%w", err)
	}

	if c.ValidateVersion {
		if err := client.CheckServerVersion(response.ServerVersion); err != nil {
			return client, err
		}
	}

	return client, nil
}

var defaultConnector AeroSpaceConnector = &AeroSpaceDefaultConnector{}

// SetDefaultConnector sets the default AeroSpaceConnector.
// This allows you to set a custom connector if needed, for testing or other purposes.
func SetDefaultConnector(connector AeroSpaceConnector) {
	if connector == nil {
		panic("ASSERTION: Default connector cannot be nil")
	}
	defaultConnector = connector
}

// GetDefaultConnector returns the default AeroSpaceConnector.
// Returns a connector to create a connection to the AeroSpace socket.
// It panics if the default connector is not initialized.
func GetDefaultConnector() AeroSpaceConnector {
	if defaultConnector == nil {
		panic("ASSERTION: Default connector is not initialized")
	}
	return defaultConnector
}
