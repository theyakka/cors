// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	// HeaderKeyReqOrigin is the http header for the request origin
	HeaderKeyReqOrigin string = "Origin"
	// HeaderKeyAccCtlReqMethod is the http header designating the CORS response allowed method
	HeaderKeyAccCtlReqMethod = "Access-Control-Request-Method"
	// HeaderKeyAccCtlReqHeaders is the http header designating the CORS response allowed headers
	HeaderKeyAccCtlReqHeaders = "Access-Control-Request-Headers"

	// HeaderKeyAccCtlResAllowOrigin is the http header designating the CORS response allowed origin
	HeaderKeyAccCtlResAllowOrigin = "Access-Control-Allow-Origin"
	// HeaderKeyAccCtlResAllowMethods is the http header designating the CORS response allowed methods
	HeaderKeyAccCtlResAllowMethods = "Access-Control-Allow-Methods"
	// HeaderKeyAccCtlResAllowHeaders is the http header designating the CORS response allowed headers
	HeaderKeyAccCtlResAllowHeaders = "Access-Control-Allow-Headers"
	// HeaderKeyAccCtlResExposeHeaders indicates which headers can be exposed as part of the response
	HeaderKeyAccCtlResExposeHeaders = "Access-Control-Expose-Headers"
	// HeaderKeyAccResCtlMaxAge is the http header designating the CORS response maximum age
	HeaderKeyAccResCtlMaxAge = "Access-Control-Max-Age"
	// HeaderKeyAccCtlResAllowCreds is the http header designating whether the CORS response allows
	// cookies / credentials
	HeaderKeyAccCtlResAllowCreds = "Access-Control-Allow-Credentials"
)

// PreflightHandlerFunc will be excuted when the preflight has completed. If it succeeds,
// the value of error will be nil. If it fails, error will contain a ValidationError that
// decribes the reason for the failure.
type PreflightHandlerFunc func(w http.ResponseWriter, r *http.Request, error *ValidationError)

// ValidatePreflight will execute the preflight flow for a request. Once the validation has
// fully executed, the handler will be executed so that you can check the response.
func (c CORS) ValidatePreflight(w http.ResponseWriter, r *http.Request, handler PreflightHandlerFunc) {
	headers := w.Header()
	// if the http method is not OPTIONS then we're going to fail because the preflight
	// should be delivered via OPTIONS. We return an error code indicating that it
	// wasn't options so that you can forward on the request if you choose.
	if r.Method != http.MethodOptions {
		handler(w, r, preflightError(PreflightErrMethodInvalid))
		return
	}

	// ensure that we don't poison any cache or force a cache to return the wrong value
	headers.Add("Vary", HeaderKeyReqOrigin)
	headers.Add("Vary", HeaderKeyAccCtlReqMethod)
	headers.Add("Vary", HeaderKeyAccCtlReqHeaders)

	// check the origin
	if c.areAllOriginsAllowed {
		// all origins are allowed, set header
		headers.Set(HeaderKeyAccCtlResAllowOrigin, "*")
	} else {
		origin := r.Header.Get(HeaderKeyReqOrigin)
		if c.IsOriginAllowed(origin) {
			// passed origin is allowed, set header
			headers.Set(HeaderKeyAccCtlResAllowOrigin, origin)
		} else {
			// the origin wasn't whitelisted
			handler(w, r, preflightError(PreflightErrOriginNotAllowed))
			return
		}
	}

	// check the requested method
	method := r.Header.Get(HeaderKeyAccCtlReqMethod)
	if method == "" {
		// the method header was missing
		handler(w, r, preflightError(PreflightErrMethodMissing))
		return
	}
	// when we compare the method we should convert to uppercase before doing the check
	upperMethod := strings.ToUpper(method)
	if c.IsMethodAllowed(upperMethod) {
		// we only return the method that was requested here.
		headers.Set(HeaderKeyAccCtlResAllowMethods, upperMethod)
	} else {
		// the method wasn't whitelisted
		handler(w, r, preflightError(PreflightErrMethodNotAllowed))
		return
	}
	// if all headers are allowed, then we should skip the check because we will need to parse the
	// header value first and that will consume time + resources
	if !c.areAllHeadersAllowed {
		// parse the header string and then check to see if the headers have been whitelisted
		headersString := r.Header.Get(HeaderKeyAccCtlReqHeaders)
		cleanedHeaders := cleanAllowedHeaderValue(headersString)
		if c.AreHeadersAllowed(cleanedHeaders) {
			headers.Set(HeaderKeyAccCtlResAllowHeaders, strings.Join(cleanedHeaders, ", "))
		} else {
			// one or more of the headers weren't whitelisted
			handler(w, r, preflightError(PreflightErrHeadersNotAllowed))
			return
		}
	}

	if len(c.exposedHeaders) > 0 {

	}

	// pass through the max age header
	if c.options.MaxAge > 0 {
		headers.Set(HeaderKeyAccResCtlMaxAge, strconv.Itoa(c.options.MaxAge))
	}

	// pass through the allow credentials header
	if c.options.AllowCredentials {
		headers.Set(HeaderKeyAccCtlResAllowCreds, "true")
	}

	handler(w, r, nil)
}

func preflightError(code int) *ValidationError {
	return preflightErrorWithSource(code, nil)
}

func preflightErrorWithSource(code int, originalError error) *ValidationError {
	message := codedErrorMessages[code]
	if message == "" {
		message = "please check code + original error for details"
	}
	return &ValidationError{
		Code:          code,
		Message:       message,
		OriginalError: originalError,
	}
}

// IsOriginAllowed does a check to see if an origin value is whitelisted according to the
// attached AllowedOrigins values.
func (c CORS) IsOriginAllowed(checkOrigin string) bool {
	// check first to see if all origins are allowed so we can get the heck out of here
	if c.areAllOriginsAllowed {
		return true
	}
	checkOrigin = strings.ToLower(checkOrigin)
	// check each of the allowed origin values to see if we have a match
	for _, origin := range c.allowedOrigins {
		if origin.AllowsFor(checkOrigin) {
			// allowed
			return true
		}
	}
	// not allowed. sorry.
	return false
}

// IsMethodAllowed will return true if the provided method value is in the list of
// whitelisted HTTP methods or it is the OPTIONS http method (which is always allowed).
func (c CORS) IsMethodAllowed(checkMethod string) bool {
	// always allow OPTIONS because it will be used for preflight
	if checkMethod == http.MethodOptions {
		return true
	}
	// check to see if the method that was passed is in the list of allowed methods
	for _, method := range c.allowedMethods {
		if method == checkMethod {
			return true
		}
	}
	// not allowed. dun dun duuunnnn.
	return false
}

func (c CORS) AreHeadersAllowed(headers []string) bool {
	for _, passedHeader := range headers {
		isAllowed := false
		for _, allowedHeader := range c.allowedHeaders {
			if passedHeader == allowedHeader {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return false
		}
	}
	// not allowed
	return true
}

func cleanAllowedHeaderValue(value string) []string {
	var headers []string
	splitValues := strings.Split(value, ",")
	for _, splitValue := range splitValues {
		trimmed := strings.TrimLeft(splitValue, " ")
		headers = append(headers, http.CanonicalHeaderKey(trimmed))
	}
	return headers
}
