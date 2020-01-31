// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

import "regexp"

// Match is a generic exact value or "wildcard" (via regex) matcher that can be used
// whenever you need to match things in the system.
type Match struct {
	Value      string
	IsWildcard bool
	regex      *regexp.Regexp
}

// NewMatch defines a new Match that should be an exact match
func NewMatch(value string) *Match {
	return &Match{
		Value:      value,
		IsWildcard: false,
		regex:      nil,
	}
}

// NewWildcardMatch defines a new Match that will have a wildcard component.
// Note: this function will automatically apply boundaries to the pattern
// to allow for exact matching of the pattern only.
func NewWildcardMatch(pattern string) *Match {
	boundaryPattern := `\b` + pattern + `\b`
	regex := regexp.MustCompile(boundaryPattern)
	return &Match{
		Value:      pattern,
		IsWildcard: true,
		regex:      regex,
	}
}

// EM is a convenience function that wraps NewMatch for Exact Matches
func EM(value string) *Match {
	return NewMatch(value)
}

// WC is a convenience function that wraps NewMatch for Wildcard Matches.
func WC(pattern string) *Match {
	match := NewWildcardMatch(pattern)
	return match
}

// Matches evaluates an input vs the Match instance to see if it matches
func (og *Match) Matches(input string) bool {
	if !og.IsWildcard {
		return og.Value == input
	}
	if og.regex == nil {
		return false
	}
	return og.regex.MatchString(input)
}
