package shared

import (
	"errors"
	"testing"
)

// TestGreet covers the normal, whitespace and blank-name cases.
func TestGreet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{name: "plain name", input: "Kiet", want: "Hello, Kiet!"},
		{name: "surrounding whitespace", input: "  Kiet \n", want: "Hello, Kiet!"},
		{name: "empty", input: "", wantErr: ErrEmptyName},
		{name: "only whitespace", input: "   ", wantErr: ErrEmptyName},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Greet(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Greet(%q) error = %v, want %v", tc.input, err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("Greet(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
