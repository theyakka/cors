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
	// HeaderKeyOrigin is the http header for the request origin
	HeaderKeyOrigin string = "Origin"
	// HeaderKeyAccCtlAllowOrigin is the http header designating the CORS response allowed origin
	HeaderKeyAccCtlAllowOrigin = "Access-Control-Allow-Origin"
	// HeaderKeyAccCtlReqMethod is the http header designating the CORS response allowed method
	HeaderKeyAccCtlReqMethod = "Access-Control-Request-Method"
	// HeaderKeyAccCtlReqHeaders is the http header designating the CORS response allowed headers
	HeaderKeyAccCtlReqHeaders = "Access-Control-Request-Headers"
	// HeaderKeyAccCtlMaxAge is the http header designating the CORS response maximum age
	HeaderKeyAccCtlMaxAge = "Access-Control-Max-Age"
	// HeaderKeyAccCtlAllowCreds is the http header designating whether the CORS response allows
	// cookies / credentials
	HeaderKeyAccCtlAllowCreds = "Access-Control-Allow-Credentials"
)

type PreflightHandler func(response PreflightResponse, w http.ResponseWriter, r *http.Request)

type PreflightResponse struct {
	Code     int
	HasError bool
}

func (c CORS) DoPreflight(w http.ResponseWriter, r *http.Request, handler PreflightHandler) {
	headers := w.Header()
	if r.Method != http.MethodOptions {
		handler(preflightError(PreflightErrMethodInvalid), w, r)
		return
	}
	// get the origin
	origin := r.Header.Get(HeaderKeyOrigin)
	// ensure that we don't poison any caches or force caches to return the wrong value
	headers.Add("Vary", HeaderKeyOrigin)
	headers.Add("Vary", HeaderKeyAccCtlReqMethod)
	headers.Add("Vary", HeaderKeyAccCtlReqHeaders)

	// check the origin
	if c.areAllOriginsAllowed {
		// all origins are allowed, set header
		headers.Set(HeaderKeyAccCtlAllowOrigin, "*")
	} else if c.IsOriginAllowed(origin) {
		// passed origin is allowed, set header
		headers.Set(HeaderKeyAccCtlAllowOrigin, origin)
	} else {
		// the origin wasn't whitelisted
		handler(preflightError(PreflightErrOriginNotAllowed), w, r)
		return
	}

	// check the requested method
	method := r.Header.Get(HeaderKeyAccCtlReqMethod)
	if method == "" {
		// the method header was missing
		handler(preflightError(PreflightErrMethodMissing), w, r)
		return
	}
	// when we compare the method we should convert to uppercase before doing the check
	upperMethod := strings.ToUpper(method)
	if c.IsMethodAllowed(upperMethod) {
		// we only return the method that was requested here.
		headers.Set(HeaderKeyAccCtlReqMethod, upperMethod)
	} else {
		// the method wasn't whitelisted
		handler(preflightError(PreflightErrMethodNotAllowed), w, r)
		return
	}
	// if all headers are allowed, then we should skip the check because we will need to parse the
	// header value first and that will consume time + resources
	if !c.areAllHeadersAllowed {
		// parse the header string and then check to see if the headers have been whitelisted
		headersString := r.Header.Get(HeaderKeyAccCtlReqHeaders)
		cleanedHeaders := cleanAllowedHeaderValue(headersString)
		if c.AreHeadersAllowed(cleanedHeaders) {
			headers.Set(HeaderKeyAccCtlReqHeaders, strings.Join(cleanedHeaders, ", "))
		} else {
			// one or more of the headers weren't whitelisted
			handler(preflightError(PreflightErrHeadersNotAllowed), w, r)
			return
		}
	}

	// pass through the max age header
	if c.options.MaxAge > 0 {
		headers.Set(HeaderKeyAccCtlMaxAge, strconv.Itoa(c.options.MaxAge))
	}

	// pass through the allow credentials header
	if c.options.AllowCredentials {
		headers.Set(HeaderKeyAccCtlAllowCreds, "true")
	}
}

func preflightError(code int) PreflightResponse {
	return PreflightResponse{Code: code, HasError: true}
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
