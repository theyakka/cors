// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

type CORS struct {
	// options is the attached set of options used to create this instance.
	options *Options
	// allowedOrigins is the cleaned list of origins we will allow. If allowedOrigins is
	// empty, then it will be assumed all origins will be allowed (areAllOriginsAllowed will
	// be true).
	allowedOrigins []*Origin
	// areAllOriginsAllowed will be true if the AllowedOrigins value in the attached Options
	// instance contained the '*' origin or if AllowedOrigins was empty.
	areAllOriginsAllowed bool
	// ...
	allowedHeaders []string
	// areAllHeadersAllowed will be true if the AllowedHeaders value in the attached Options
	// instance contained the '*' origin. areAllHeadersAllowed will NOT be true if
	// AllowedHeaders are empty as we will use a default set of headers. See docs for
	// AllowedHeaders for details.
	areAllHeadersAllowed bool
	// allowedMethods is a cleaned list of all of the HTTP methods that will be allowed.
	allowedMethods []string
}

// NewCORS creates a new CORS instance that is pre-configured with the values
// defined in the provided Options instance.
func NewCORS(options Options) (*CORS, error) {
	return options.NewCORS()
}

// AllowAll creates a new CORS instance that allows all origins, methods, and
// headers.
func AllowAll() (*CORS, error) {
	return NewCORS(Options{
		AllowedOrigins: AllowAllOrigins,
		AllowedMethods: AllowAllMethods,
		AllowedHeaders: AllowAllHeaders,
	})
}
