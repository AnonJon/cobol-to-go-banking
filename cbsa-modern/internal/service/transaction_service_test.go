package service

import "testing"

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"deadlock error", errStr("pq: deadlock detected"), true},
		{"serialization error", errStr("could not serialize access"), true},
		{"normal error", errStr("account not found"), false},
		{"timeout with deadlock", errStr("some deadlock happened"), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isRetryable(tc.err)
			if got != tc.expected {
				t.Errorf("isRetryable(%v) = %v, want %v", tc.err, got, tc.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s, substr string
		expected  bool
	}{
		{"hello world", "world", true},
		{"hello", "world", false},
		{"deadlock detected", "deadlock", true},
		{"", "a", false},
		{"a", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.s+"/"+tc.substr, func(t *testing.T) {
			got := contains(tc.s, tc.substr)
			if got != tc.expected {
				t.Errorf("contains(%q, %q) = %v, want %v", tc.s, tc.substr, got, tc.expected)
			}
		})
	}
}

type errStr string

func (e errStr) Error() string { return string(e) }
