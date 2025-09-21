package webhook

import (
	"testing"
)

// TestIsSupportedEvent tests the webhook event filtering
func TestIsSupportedEvent(t *testing.T) {
	tests := []struct {
		eventType string
		expected  bool
	}{
		{"push", true},
		{"issue_comment", true},
		{"pull_request", true},
		{"ping", false},
		{"release", false},
		{"issues", false},
		{"fork", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsSupportedEvent(test.eventType)
		if result != test.expected {
			t.Errorf("IsSupportedEvent(%q) = %v, expected %v", test.eventType, result, test.expected)
		}
	}
}

// TestSupportedEventTypes verifies the supported event types map
func TestSupportedEventTypes(t *testing.T) {
	expected := map[string]bool{
		"push":          true,
		"issue_comment": true,
		"pull_request":  true,
	}

	for eventType, expectedValue := range expected {
		if SupportedEventTypes[eventType] != expectedValue {
			t.Errorf("SupportedEventTypes[%q] = %v, expected %v", eventType, SupportedEventTypes[eventType], expectedValue)
		}
	}

	// Test that unsupported events are not in the map (or false)
	unsupportedEvents := []string{"ping", "release", "issues", "fork"}
	for _, eventType := range unsupportedEvents {
		if SupportedEventTypes[eventType] {
			t.Errorf("SupportedEventTypes[%q] should be false or not present", eventType)
		}
	}
}