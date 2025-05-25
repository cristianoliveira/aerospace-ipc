package aerospace

import "testing"

func TestWindows(t *testing.T) {
	t.Run("formatting window as string", func(tt *testing.T) {
		testCases := []struct {
			title    string
			window  Window
			expected string
		}{
			{
				title:   "Basic Window Formatting",
				window:  Window{WindowID: 123, WindowTitle: "Test Window", AppName: "TestApp"},
				expected: "123 | TestApp | Test Window",
			},
			{
				title:   "Basic 2 Window Formatting",
				window:  Window{WindowID: 456, WindowTitle: "Another Window", AppName: "AnotherApp"},
				expected: "456 | AnotherApp | Another Window",
			},
			{
				title:   "Window with Empty App Name",
				window:  Window{WindowID: 789, WindowTitle: "Sample Window", AppName: ""},
				expected: "789 |  | Sample Window",
			},
			{
				title:   "Window with Empty Title",
				window:  Window{WindowID: 101, WindowTitle: "", AppName: "EmptyTitleApp"},
				expected: "101 | EmptyTitleApp",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.title, func(t *testing.T) {
				result := tc.window.String()
				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			})
		}
	})
}
