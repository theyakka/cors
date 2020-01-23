// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

// Options represents the configurable elements of the CORS validation
// process.
type Options struct {
	// The list of origins you want to whitelist
	AllowedOrigins []string
	// The list of methods you want to whitelist
	AllowedMethods []string
	// The list of headers you want to whitelist
	AllowedHeaders []string
	MaxAge         int
	// AllowCredentials, when set to true, will allow the request to include
	// credentials such as cookies or otherwise. Note, you cannot set the value
	// to true AND use wildcard values for other Options values. If you attempt
	// to do so, there will be a configuration error. The default value is false.
	AllowCredentials bool
}

// OptionsAllowAll creates a default set of options that allows all origins,
// all common HTTP methods and a common set of headers.
func OptionsAllowAll() *Options {
	return &Options{
		AllowedOrigins:   AllowAllOrigins,
		AllowedMethods:   AllowAllMethods,
		AllowedHeaders:   DefaultAllowedHeaders,
		AllowCredentials: false,
		MaxAge:           0,
	}
}

// NewCORS creates a new CORS instance that is pre-configured with the values
// defined in the current Options instance.
func (o *Options) NewCORS() (*CORS, error) {
	c := &CORS{}
	o.applyAllowedOrigins(c)
	o.applyAllowedMethods(c)
	o.applyAllowedHeaders(c)
	if o.AllowCredentials && (c.areAllOriginsAllowed || c.areAllHeadersAllowed) {
		return nil, errors.New("")
	}
	c.options = o
	return c, nil
}

// applyAllowedOrigins checks for the "*" value and will
func (o *Options) applyAllowedOrigins(c *CORS) {
	var allowedOrigins []*Origin
	if len(o.AllowedOrigins) == 0 {
		// no allowed origins so we're going to assume that you want to allow everything
		// (vs nothing .. as that would be weird)
		c.areAllOriginsAllowed = true
		c.allowedOrigins = nil
	} else {
		for _, originString := range o.AllowedOrigins {
			if originString == "*" {
				// if we're allowing all origins then it's irrelevant to keep the old list
				// so we will just stop here
				c.areAllOriginsAllowed = true
				c.allowedOrigins = nil
				break
			}
			origin, originErr := NewOrigin(strings.ToLower(originString))
			if originErr != nil {
				log.Println(originErr)
			}
			allowedOrigins = append(allowedOrigins, origin)
		}
		c.areAllOriginsAllowed = false
		c.allowedOrigins = allowedOrigins
	}
}

func (o *Options) applyAllowedMethods(c *CORS) {
	if len(o.AllowedMethods) == 0 {
		// use the simple request HTTP method types that the spec defines as the default
		// because nothing was passed.
		c.allowedMethods = SpecSimpleMethods
	} else {
		// convert any provided value to uppercase so that, later, when we do our checks
		// we don't waste time converting them every time.
		var allowedMethods []string
		for _, method := range o.AllowedMethods {
			allowedMethods = append(allowedMethods, strings.ToUpper(method))
		}
		c.allowedMethods = allowedMethods
	}
}

func (o *Options) applyAllowedHeaders(c *CORS) {
	if len(o.AllowedHeaders) == 0 {
		c.areAllHeadersAllowed = false
		c.allowedHeaders = DefaultAllowedHeaders
	} else {
		var allowedHeaders []string
		for _, header := range o.AllowedHeaders {
			if header == "*" {
				c.areAllHeadersAllowed = true
				c.allowedHeaders = nil
				break
			}
			allowedHeaders = append(allowedHeaders, http.CanonicalHeaderKey(header))
		}
		c.areAllHeadersAllowed = false
		c.allowedHeaders = allowedHeaders
	}

}

// DefaultHeadersWith returns a slice containing the headers you provided, along
// with the default set of headers. It does not check for duplicate entries. See
// DefaultAllowedHeaders for default headers value.
func DefaultHeadersWith(headers ...string) []string {
	return append(DefaultAllowedHeaders, headers...)
}

// AllowAllOrigins is a slice containing just the wildcard ("*") origin.
var AllowAllOrigins = []string{"*"}

// DefaultAllowedHeaders is a list of the common headers that you will want to allow for
// all CORS preflights / requests. It is used as the default list if you don't specify
// anything.
var DefaultAllowedHeaders = []string{
	"Accept", "Content-Type", "Origin", "X-Requested-With",
}

// SpecSimpleMethods contains the HTTP methods that the CORS specification deems acceptable
// methods for "simple" requests. We use these methods as the default if value is provided
// for AllowedMethods.
var SpecSimpleMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
}

// AllowAllMethods is a list of all common HTTP methods.
var AllowAllMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
}

// AllowAllHeaders is a slice containing just the wildcard ("*") header.
var AllowAllHeaders = []string{"*"}

// DefaultExposeHeaders
var DefaultExposeHeaders = []string{
	"Cache-Control", "Content-Language", "Content-Type",
	"Expires", "Last-Modified", "Pragma",
}
