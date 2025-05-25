package aerospace

import (
	"fmt"
	"net"
)

// Connector should return AeroSpaceConnectiono
// AeroSpaceConnector is an interface for connecting to the AeroSpace socket.
//
// It provides a method to establish a connection and return an AeroSpaceSocketConn.
//  See: AeroSpaceDefaultConnector for the default implementation.
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
	socketPath, err := GetSocketPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get socket path\n %w", err)
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceSocketConnection{
		MinAerospaceVersion: "0.15.2-Beta",
		Conn:                &conn,
		socketPath:          socketPath,
	}

	return client, nil
}

var DefaultConnector AeroSpaceConnector = &AeroSpaceDefaultConnector{}
