package aerospace

import (
	"fmt"
	"net"

	"github.com/cristianoliveira/aerospace-ipc/client"
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
	Connect() (client.AeroSpaceSocketConn, error)
}

// AeroSpaceDefaultConnector is the default implementation of AeroSpaceConnector.
//
// In most cases, you will use this connector to connect to the AeroSpace socket.
type AeroSpaceDefaultConnector struct{}

func (c *AeroSpaceDefaultConnector) Connect() (client.AeroSpaceSocketConn, error) {
	socketPath, err := GetSocketPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get socket path\n %w", err)
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &client.AeroSpaceSocketConnection{
		MinAerospaceVersion: AeroSpaceSocketClientVersion,
		Conn:                &conn,
		SocketPath:          socketPath,
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
func GetDefaultConnector() AeroSpaceConnector {
	if defaultConnector == nil {
		panic("ASSERTION: Default connector is not initialized")
	}
	return defaultConnector
}
