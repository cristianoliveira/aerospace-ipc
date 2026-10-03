package monitors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// Monitor describes a monitor managed by AeroSpaceWM.
type Monitor struct {
	// MonitorID is the 1-based ID of the monitor.
	MonitorID int `json:"monitor-id"`

	// MonitorName is the monitor's name.
	MonitorName string `json:"monitor-name"`
}

// MonitorsService defines the monitor query operations in AeroSpaceWM.
type MonitorsService interface {
	// GetFocusedMonitor returns the currently focused monitor.
	GetFocusedMonitor() (*Monitor, error)

	// GetAllMonitors returns all monitors.
	GetAllMonitors() ([]Monitor, error)
}

// Service provides methods to query monitors in AeroSpaceWM.
type Service struct {
	client client.AeroSpaceConnection
}

// NewService creates a monitor service with the given AeroSpace connection.
func NewService(client client.AeroSpaceConnection) *Service {
	return &Service{client: client}
}

// GetFocusedMonitor returns the currently focused monitor.
//
// It is equivalent to running the command:
//
//	aerospace list-monitors --focused --json --format "%{monitor-id} %{monitor-name}"
//
// Returns an error if no monitor is focused, the command fails, or its output is malformed.
func (s *Service) GetFocusedMonitor() (*Monitor, error) {
	response, err := s.client.SendCommand(
		"list-monitors",
		[]string{
			"--focused",
			"--json",
			"--format",
			"%{monitor-id} %{monitor-name}",
		},
	)
	if err != nil {
		return nil, err
	}
	if response.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get focused monitor:\n%s", response.StdErr)
	}

	monitors, err := parseMonitorList(response.StdOut)
	if err != nil {
		return nil, fmt.Errorf("failed to parse focused monitor: %w", err)
	}
	if len(monitors) == 0 {
		return nil, fmt.Errorf("no focused monitor found")
	}
	if len(monitors) > 1 {
		return nil, fmt.Errorf("expected one focused monitor, got %d", len(monitors))
	}

	return &monitors[0], nil
}

// GetAllMonitors returns all monitors in AeroSpace's current order.
//
// It is equivalent to running the command:
//
//	aerospace list-monitors --json --format "%{monitor-id} %{monitor-name}"
//
// An empty array is a valid empty inventory. Malformed or incomplete monitor data is rejected.
func (s *Service) GetAllMonitors() ([]Monitor, error) {
	response, err := s.client.SendCommand(
		"list-monitors",
		[]string{
			"--json",
			"--format",
			"%{monitor-id} %{monitor-name}",
		},
	)
	if err != nil {
		return nil, err
	}
	if response.ExitCode != 0 {
		return nil, fmt.Errorf("failed to list monitors:\n%s", response.StdErr)
	}

	monitors, err := parseMonitorList(response.StdOut)
	if err != nil {
		return nil, fmt.Errorf("failed to parse monitors: %w", err)
	}
	return monitors, nil
}

func parseMonitorList(output string) ([]Monitor, error) {
	var monitors []Monitor
	if err := json.Unmarshal([]byte(output), &monitors); err != nil {
		return nil, err
	}
	if monitors == nil {
		return nil, fmt.Errorf("expected a JSON array")
	}

	seenIDs := make(map[int]struct{}, len(monitors))
	for i, monitor := range monitors {
		if monitor.MonitorID <= 0 {
			return nil, fmt.Errorf("monitor at index %d has nonpositive ID %d", i, monitor.MonitorID)
		}
		if strings.TrimSpace(monitor.MonitorName) == "" {
			return nil, fmt.Errorf("monitor at index %d has a blank name", i)
		}
		if _, exists := seenIDs[monitor.MonitorID]; exists {
			return nil, fmt.Errorf("duplicate monitor ID %d", monitor.MonitorID)
		}
		seenIDs[monitor.MonitorID] = struct{}{}
	}

	return monitors, nil
}
