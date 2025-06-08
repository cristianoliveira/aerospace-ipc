package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
)

// Command represents the JSON structure for AeroSpace socket commands.
type Command struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Stdin   string   `json:"stdin"`
}

// Response represents the JSON structure from AeroSpace socket response.
type Response struct {
	ServerVersion string `json:"serverVersionAndHash"` // Fornat: "0.0.1-Beta <hash>"
	StdErr        string `json:"stderr"`
	StdOut        string `json:"stdout"`
	ExitCode      int32  `json:"exitCode"`
}

// AeroSpaceConnection is an interface interacting with a AeroSpace socket.
//
// It provides methods to execute low-level commands and manage the connection.
type AeroSpaceConnection interface {
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

	// GetServerVersion returns the version of the AeroSpace server.
	GetServerVersion() (string, error)

	// CheckServerVersion validates the version of the AeroSpace server.
	CheckServerVersion() error
}

// AeroSpaceSocketConnection implements the AeroSpaceSocketConn interface.
type AeroSpaceSocketConnection struct {
	mu              sync.Mutex
	socketPath      string
	MinMajorVersion int
	MinMinorVersion int
	Conn            net.Conn
}

// GetSocketPath returns the socket path for the AeroSpace connection.
//
// It returns an error if the socket path is not set.
func (c *AeroSpaceSocketConnection) GetSocketPath() (string, error) {
	if c.socketPath == "" {
		return "", fmt.Errorf("missing socket path")
	}

	return c.socketPath, nil
}

// CloseConnection closes the connection to the AeroSpace socket.
//
// It returns an error if the connection cannot be closed.
func (c *AeroSpaceSocketConnection) CloseConnection() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection\n %w", err)
		}
	}
	return nil
}

// GetServerVersion retrieves the version of the AeroSpace server.
// It sends a command to the server to get its version and returns it.
func (c *AeroSpaceSocketConnection) GetServerVersion() (string, error) {
	if c.Conn == nil {
		return "", fmt.Errorf("connection is not established")
	}
	// Send the command to get the server version
	res, err := c.SendCommand("config", []string{"--config-path"})
	if err != nil {
		return "", fmt.Errorf("failed to get server version\n%w", err)
	}

	if res.ExitCode != 0 {
		return "", fmt.Errorf("failed to get server version\n%s", res.StdErr)
	}

	return res.ServerVersion, nil
}

// CheckServerVersion checks if the server version meets the minimum requirements.
// It compares the server version against the minimum major and minor versions.
func (c *AeroSpaceSocketConnection) CheckServerVersion() error {
	serverVersion, err := c.GetServerVersion()
	if err != nil {
		return fmt.Errorf("failed to get server version\n%w", err)
	}

	if serverVersion == "" {
		return fmt.Errorf("server version is empty")
	}
	parts := strings.Split(serverVersion, "-")
	versionParts := strings.Split(parts[0], ".")
	if len(versionParts) < 2 {
		fmt.Printf("[WARN] Invalid server version format: %s\n", serverVersion)
	}

	intMajor, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("failed to parse major version from %s\n%w", serverVersion, err)
	}

	intMinor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return fmt.Errorf("failed to parse minor version from %s\n%w", serverVersion, err)
	}

	if intMajor != c.MinMajorVersion ||
		intMajor == c.MinMajorVersion && intMinor != c.MinMinorVersion {
		versionJoined := strings.Join(versionParts, ".")
		return exceptions.NewErrVersionMismatch(
			c.MinMajorVersion,
			c.MinMinorVersion,
			versionJoined,
		)
	}

	return nil
}

// SendCommand sends a raw command to the AeroSpace socket and returns a raw response.
// It allows to execute commands that are not directly supported by the client library.
//
// It is equivalent to running the command:
//
//	aerospace <command> <args...>
//
// Returns a Response struct containing the server version, standard error, standard output, and exit code.
//
// Usage:
//
//	response, err := client.SendCommand("list-windows", []string{"--all", "--json"})
//	if err != nil {
//	  fmt.Println("Error:", err)
//	}
//	fmt.Println("Server Version:", response.ServerVersion)
//	fmt.Println("Standard Output:", response.StdOut)
//	fmt.Println("Standard Error:", response.StdErr)
func (c *AeroSpaceSocketConnection) SendCommand(command string, args []string) (*Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
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

	_, err = c.Conn.Write(cmdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send command\n%w", err)
	}

	var responseData []byte
	buf := make([]byte, 4096)
	for {
		n, err := c.Conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf(
				"failed to read response\n%w\ndata\n%s",
				err,
				responseData,
			)
		}
		responseData = append(responseData, buf[:n]...)
		if n < len(buf) {
			break
		}
	}

	var response Response
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal socket response\n%w\ndata\n%s",
			err,
			responseData,
		)
	}

	if response.ExitCode != 0 {
		return nil, fmt.Errorf("command failed with exit code %d\n%s", response.ExitCode, response.StdErr)
	}

	if response.StdErr != "" {
		return nil, fmt.Errorf("command error\n%s", response.StdErr)
	}

	return &response, nil
}

// NewAeroSpaceSocketConnection creates a new AeroSpaceSocketConnection.
// It initializes the connection to the AeroSpace socket.
func NewAeroSpaceSocketConnection(socketPath string) (*AeroSpaceSocketConnection, error) {
	if socketPath == "" {
		return nil, fmt.Errorf("socket path cannot be empty")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceSocketConnection{
		socketPath:      socketPath,
		MinMajorVersion: 1, // Example version, adjust as needed
		MinMinorVersion: 0,
		Conn:            conn,
	}

	return client, nil
}
