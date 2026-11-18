// Package shared holds code used by more than one service. It has no image of its own: CI lints
// and tests it through the top-level x-ci.checks entry, and rebuilds api when it changes (api's
// x-ci.watch).
package shared

import (
	"errors"
	"strings"
)

// ErrEmptyName is returned by Greet when the name is blank.
var ErrEmptyName = errors.New("name must not be empty")

// Greet builds the greeting returned by the api service.
//
// Parameters:
//   - name: the person to greet; surrounding whitespace is ignored.
//
// Returns the greeting, or ErrEmptyName when name is blank.
func Greet(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", ErrEmptyName
	}
	return "Hello, " + trimmed + "!", nil
}
