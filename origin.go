// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

import "strings"

import "regexp"

type Origin struct {
	Value      string
	IsWildcard bool
	regex      *regexp.Regexp
}

func NewOrigin(origin string) (*Origin, error) {
	isWildcard := strings.Contains(origin, "*")
	var regex *regexp.Regexp
	var regexErr error
	if isWildcard {
		regOrigin := strings.ReplaceAll(origin, ".", "\\.")
		regOrigin = strings.ReplaceAll(regOrigin, "*", ".*")
		regex, regexErr = regexp.Compile(regOrigin)
		if regexErr != nil {
			return nil, regexErr
		}
	}
	return &Origin{
		Value:      origin,
		IsWildcard: isWildcard,
		regex:      regex,
	}, nil
}

func (og Origin) AllowsFor(origin string) bool {
	if !og.IsWildcard {
		return og.Value == origin
	}
	return og.regex.MatchString(origin)
}
