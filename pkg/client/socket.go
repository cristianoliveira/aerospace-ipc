package client

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cristianoliveira/aerospace-ipc/internal/constants"
	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
)

// socketProtocolVersion must match SOCKET_PROTOCOL_VERSION in AeroSpace's clientServer.swift.
// AeroSpace added this handshake in commit 8413641c (2026-06-21).
const socketProtocolVersion uint32 = 1

// Command represents the JSON structure for AeroSpace socket commands.
// This wlll mostly mirror https://github.com/nikitabobko/AeroSpace/blob/main/Sources/Common/model/clientServer.swift#L76
type Command struct {
	Args  []string `json:"args"`
	Stdin string   `json:"stdin"`
	// Pass null if callback context is unavailable.
	WindowID  *uint32 `json:"windowId"`
	Workspace *string `json:"workspace"`
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
		return fmt.Errorf("invalid server version format: %s", serverVersion)
	}

	intMajor, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("failed to parse major version from %s\n%w", serverVersion, err)
	}

	intMinor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return fmt.Errorf("failed to parse minor version from %s\n%w", serverVersion, err)
	}

	// Since AeroSpace may have breaking changes even in minor versions,
	// I'll enforce exact match on major and minor versions
	// Socket protocol version 1 requires AeroSpace 0.21.0 or newer.
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
	if command == "subscribe" {
		return nil, fmt.Errorf("subscribe is not supported by SendCommand; use a dedicated streaming connection")
	}

	commandArgs := append([]string{command}, args...)
	cmd := Command{
		Args:  commandArgs,
		Stdin: "",
	}

	windowID, ok := os.LookupEnv("AEROSPACE_WINDOW_ID")
	if ok {
		parsedID, err := strconv.ParseUint(windowID, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AEROSPACE_WINDOW_ID\n%w", err)
		}
		windowID := uint32(parsedID)
		cmd.WindowID = &windowID
	}

	// Pass env AEROSPACE_WORKSPACE if available
	workspace, ok := os.LookupEnv("AEROSPACE_WORKSPACE")
	if ok {
		cmd.Workspace = &workspace
	}

	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command\n%w", err)
	}

	// AeroSpace expects [4-byte little-endian length][JSON payload].
	frame := make([]byte, 4, 4+len(cmdBytes))
	binary.LittleEndian.PutUint32(frame, uint32(len(cmdBytes)))
	frame = append(frame, cmdBytes...)
	if err = writeFull(c.Conn, frame); err != nil {
		return nil, fmt.Errorf("failed to send command\n%w", err)
	}

	// Response is also length-prefixed
	var responseLen uint32
	if err := binary.Read(c.Conn, binary.LittleEndian, &responseLen); err != nil {
		return nil, fmt.Errorf("failed to read response length\n%w", err)
	}
	const maxResponseSize = 16 * 1024 * 1024 // 16 MiB
	if responseLen > maxResponseSize {
		return nil, fmt.Errorf("response length %d exceeds maximum %d", responseLen, maxResponseSize)
	}
	responseData := make([]byte, responseLen)
	if _, err := io.ReadFull(c.Conn, responseData); err != nil {
		return nil, fmt.Errorf("failed to read response\n%w", err)
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

	return &response, nil
}

func writeFull(writer io.Writer, data []byte) error {
	for len(data) > 0 {
		written, err := writer.Write(data)
		if err != nil {
			return err
		}
		if written <= 0 || written > len(data) {
			return io.ErrShortWrite
		}
		data = data[written:]
	}
	return nil
}

func closeAfterConnectionError(conn net.Conn, connectionErr error) error {
	if closeErr := conn.Close(); closeErr != nil {
		return errors.Join(connectionErr, fmt.Errorf("failed to close connection: %w", closeErr))
	}
	return connectionErr
}

// NewAeroSpaceSocketConnection creates a new AeroSpaceSocketConnection.
// It initializes the connection to the AeroSpace socket and performs the protocol version handshake.
func NewAeroSpaceSocketConnection(socketPath string) (*AeroSpaceSocketConnection, error) {
	if socketPath == "" {
		return nil, fmt.Errorf("socket path cannot be empty")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	// AeroSpace requires a protocol version handshake on connect (added in commit 8413641c).
	// Client sends its version; server responds with its own; both must match.
	if err := binary.Write(conn, binary.LittleEndian, socketProtocolVersion); err != nil {
		connectionErr := fmt.Errorf("failed to send protocol version\n%w", err)
		return nil, closeAfterConnectionError(conn, connectionErr)
	}
	var serverVersion uint32
	if err := binary.Read(conn, binary.LittleEndian, &serverVersion); err != nil {
		connectionErr := fmt.Errorf("failed to read server protocol version\n%w", err)
		return nil, closeAfterConnectionError(conn, connectionErr)
	}
	if serverVersion != socketProtocolVersion {
		connectionErr := fmt.Errorf(
			"socket protocol version mismatch: client=%d server=%d (try restarting AeroSpace)",
			socketProtocolVersion, serverVersion,
		)
		return nil, closeAfterConnectionError(conn, connectionErr)
	}

	client := &AeroSpaceSocketConnection{
		socketPath:      socketPath,
		MinMajorVersion: constants.AeroSpaceSocketClientMajor,
		MinMinorVersion: constants.AeroSpaceSocketClientMinor,
		Conn:            conn,
	}

	return client, nil
}
