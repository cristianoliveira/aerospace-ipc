package aerospace

import (
	"encoding/json"
	"fmt"
	"net"
)

// Connector should return AeroSpaceConnectiono

// Command represents the JSON structure for AeroSpace socket commands.
type Command struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Stdin   string   `json:"stdin"`
}

// Response represents the JSON structure from AeroSpace socket response.
type Response struct {
	ServerVersion string `json:"serverVersion"`
	StdErr        string `json:"stderr"`
	StdOut        string `json:"stdout"`
	ExitCode      int32  `json:"exitCode"`
}

// AeroSpaceSocketConn is an interface interacting with a AeroSpace socket.
// 
// It provides methos to execute low-level commands and manage the connection.
type AeroSpaceSocketConn interface {
	// CloseConnection closes the connection to the AeroSpace socket.
	CloseConnection() error

	// SendCommand sends a raw command to the AeroSpace socket and returns a raw response.
	// 
	// It is equivalent to running the command:
	//   aerospace <command> <args...>
	//
	// Returns a Response struct containing the server version, standard error, standard output, and exit code.
	SendCommand(command string, args []string) (*Response, error)

	// GetSocketPath returns the socket path for the AeroSpace connection.
	GetSocketPath() (string, error)
}

// AeroSpaceSocketConnection implements the AeroSpaceSocketConn interface.
type AeroSpaceSocketConnection struct {
	SocketPath          string
	MinAerospaceVersion string
	Conn                *net.Conn
	socketPath          string
}

func (c *AeroSpaceSocketConnection) GetSocketPath() (string, error) {
	if c.SocketPath != "" {
		return "", fmt.Errorf("missing socket path")
	}

	return c.socketPath, nil
}

func (c *AeroSpaceSocketConnection) CloseConnection() error {
	if c.Conn != nil {
		err := (*c.Conn).Close()
		if err != nil {
			return fmt.Errorf("failed to close connection\n %w", err)
		}
	}
	return nil
}

// SendCommand sends a command to the AeroSpace window manager via Unix socket and returns the response.
func (c *AeroSpaceSocketConnection) SendCommand(command string, args []string) (*Response, error) {
	if c.Conn == nil {
		return nil, fmt.Errorf("connection is not established")
	}

	// Merge command and arguments into the Command struct
	commandArgs := append([]string{command}, args...)
	cmd := Command{
		Command: "", // This field is deprecated and not used
		Args:    commandArgs,
		Stdin:   "",
	}

	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command\n%w", err)
	}

	_, err = (*c.Conn).Write(cmdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send command\n%w", err)
	}

	buf := make([]byte, 4096)
	n, err := (*c.Conn).Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response\n%w", err)
	}

	var response Response
	err = json.Unmarshal(buf[:n], &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response\n%w", err)
	}

	if response.ExitCode != 0 {
		return nil, fmt.Errorf("command failed with exit code %d\n%s", response.ExitCode, response.StdErr)
	}

	if response.StdErr != "" {
		return nil, fmt.Errorf("command error\n %s", response.StdErr)
	}

	return &response, nil
}

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
