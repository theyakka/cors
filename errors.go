// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

import "fmt"

const (
	// ConfigurationInvalid means that you tried to configure the CORS instance or
	// defined an Options combination that was invalid.
	ConfigurationInvalid int = iota + 100
	// PreflightErrOriginNotAllowed means that the preflight failed because the origin
	// was not whitelisted.
	PreflightErrOriginNotAllowed
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

// codedErrorMessages is a map of user friendly error messages for the numeric error
// codes used in the system.
var codedErrorMessages = map[int]string{
	ConfigurationInvalid:          "one or more options were invalid",
	PreflightErrOriginNotAllowed:  "the requested origin was not whitelisted",
	PreflightErrMethodNotAllowed:  "the requested method was not whitelisted",
	PreflightErrHeadersNotAllowed: "one or more headers were not whitelisted",
	PreflightErrMethodMissing:     "you did not provide a http method for validation",
	PreflightErrMethodInvalid:     "you atempted to validate a CORS request but you did the request was not sent using the OPTIONS http method",
}

// ValidationError will be thrown whenever there are validation or configuration issues
type ValidationError struct {
	// Code provides a code that indicates the specific error condition
	Code int `json:"code"`
	// Message is a human readable explanation of the error condition
	Message string `json:"message"`
	// OriginalError contains the underlying source error (if any exists)
	OriginalError error `json:"error,omitempty"`
}

// Error implements the builtin error interface for our custom ValidationError type
func (ce ValidationError) Error() string {
	return fmt.Sprintf("%s [%d]", ce.Message, ce.Code)
}
