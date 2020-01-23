// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

import "fmt"

const (
	// PreflightErrOriginNotAllowed means that the preflight failed because the origin
	// was not whitelisted.
	PreflightErrOriginNotAllowed int = iota + 100
	// PreflightErrMethodNotAllowed means that the preflight failed because the http
	// method was not whitelisted.
	PreflightErrMethodNotAllowed
	// PreflightErrHeadersNotAllowed means that the preflight failed because one or more
	// headers were not whitelisted.
	PreflightErrHeadersNotAllowed
	// PreflightErrMethodMissing means you're hitting the preflight but you the
	// Access-Control-Request-Method header was not passed.
	PreflightErrMethodMissing
	// PreflightErrMethodInvalid means you're hitting the preflight but you aren't
	// using the OPTIONS method.
	PreflightErrMethodInvalid
)

// ValidationError will be thrown whenever there are validation or configuration issues
type ValidationError struct {
	Code          int
	Message       string
	OriginalError error
}

func (ce ValidationError) Error() string {
	return fmt.Sprintf("%s [%d]", ce.Message, ce.Code)
}
