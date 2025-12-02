package layout

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// SetLayoutOpts contains optional parameters for SetLayout.
type SetLayoutOpts struct {
	// WindowID specifies the window ID to set layout for. If not set, the focused window is used.
	WindowID *int
}

// Service provides methods to interact with layout in AeroSpaceWM.
type Service struct {
	client client.AeroSpaceConnection
}

// LayoutService defines the interface for layout operations in AeroSpaceWM.
type LayoutService interface {
	// SetLayout sets the layout for the focused window or a specific window.
	SetLayout(layouts []string, opts ...SetLayoutOpts) error
}

// NewService creates a new layout service with the given AeroSpace client connection.
func NewService(client client.AeroSpaceConnection) *Service {
	return &Service{client: client}
}

// SetLayout sets the layout for the focused window or a specific window.
//
// Layouts can be one or more of: accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
// If multiple layouts are provided, finds the first that doesn't describe the currently active layout and applies it.
// This is useful for toggling between layouts.
//
// It is equivalent to running the command:
//
//	aerospace layout <layout>... [--window-id <window-id>]
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Set a single layout for focused window
//	err := layoutService.SetLayout([]string{"floating"})
//
//	// Toggle between layouts (order doesn't matter)
//	err := layoutService.SetLayout([]string{"floating", "tiling"})
//	err := layoutService.SetLayout([]string{"horizontal", "vertical"})
//
//	// Set layout for specific window
//	err := layoutService.SetLayout([]string{"floating"}, layout.SetLayoutOpts{
//	    WindowID: layout.IntPtr(12345),
//	})
//
//	// Toggle layout for specific window
//	err := layoutService.SetLayout([]string{"floating", "tiling"}, layout.SetLayoutOpts{
//	    WindowID: layout.IntPtr(12345),
//	})
func (s *Service) SetLayout(layouts []string, opts ...SetLayoutOpts) error {
	if len(layouts) == 0 {
		return fmt.Errorf("at least one layout must be provided")
	}

	cmdArgs := make([]string, 0, len(layouts)+2)
	cmdArgs = append(cmdArgs, layouts...)

	var opt SetLayoutOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.WindowID != nil {
		cmdArgs = append(cmdArgs, "--window-id", fmt.Sprintf("%d", *opt.WindowID))
	}

	response, err := s.client.SendCommand("layout", cmdArgs)
	if err != nil {
		return fmt.Errorf("failed to set layout(s) %v: %w", layouts, err)
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to set layout(s) %v\n%s", layouts, response.StdErr)
	}

	return nil
}

// Helper functions for creating pointers (useful for API usage)

// IntPtr returns a pointer to the given int value.
func IntPtr(v int) *int {
	return &v
}
