// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors

// CORS is the primary manager for all the CORS validation / utility functions in this
// library. You will use an instance of this type for all interactions. Instances can be
// created directly or they can be created via Options.NewCORS.
type CORS struct {
	// options is the attached set of options used to create this instance.
	options *Options
	// allowedOrigins is the cleaned list of origins we will allow. If allowedOrigins is
	// empty, then it will be assumed all origins will be allowed (areAllOriginsAllowed will
	// be true).
	allowedOrigins []*Match
	// areAllOriginsAllowed will be true if the AllowedOrigins value in the attached Options
	// instance contained the '*' origin or if AllowedOrigins was empty.
	areAllOriginsAllowed bool
	// allowedMethods is a cleaned list of all of the HTTP methods that will be allowed.
	allowedMethods []string
	// allowedHeader is the cleaned list of all of the headers we will allow. If empty, and
	// areAllHeadersAllowed is false, then no headers will be allowed.
	allowedHeaders []string
	// areAllHeadersAllowed will be true if the AllowedHeaders value in the attached Options
	// instance contained the '*' origin. areAllHeadersAllowed will NOT be true if
	// AllowedHeaders are empty as we will use a default set of headers. See docs for
	// AllowedHeaders for details.
	areAllHeadersAllowed bool
	//
	exposedHeaders []string
}

// AllowAll creates a new CORS instance that allows all origins, methods, and
// headers.
func AllowAll() (*CORS, error) {
	return OptionsAllowAll().NewCORS()
}
