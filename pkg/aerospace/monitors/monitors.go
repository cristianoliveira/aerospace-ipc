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
// \taerospace list-monitors --focused --json --format "%{monitor-id} %{monitor-name}"
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

	var monitors []Monitor
	if err := json.Unmarshal([]byte(response.StdOut), &monitors); err != nil {
		return nil, fmt.Errorf("failed to parse focused monitor: %w", err)
	}
	if monitors == nil {
		return nil, fmt.Errorf("failed to parse focused monitor: expected a JSON array")
	}
	if len(monitors) == 0 {
		return nil, fmt.Errorf("no focused monitor found")
	}
	if len(monitors) > 1 {
		return nil, fmt.Errorf("expected one focused monitor, got %d", len(monitors))
	}

	monitor := monitors[0]
	if monitor.MonitorID <= 0 {
		return nil, fmt.Errorf("invalid focused monitor: monitor ID must be positive")
	}
	if strings.TrimSpace(monitor.MonitorName) == "" {
		return nil, fmt.Errorf("invalid focused monitor: monitor name is blank")
	}

	return &monitor, nil
}
