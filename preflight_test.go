// Created by Yakka (https://theyakka.com)
//
// Copyright (c) 2020 Yakka LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package cors_test

import (
	"github.com/theyakka/cors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSRequest(t *testing.T) {
	o := cors.Options{
		AllowedOrigins: []string{"http*://*.theyakka.com", "http*://theyakka.com"},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		t.Error(err)
		return
	}

	url := "https://theyakka.com"
	req, _ := http.NewRequest(http.MethodOptions, url, nil)
	req.Header.Set(cors.HeaderKeyOrigin, url)
	req.Header.Set(cors.HeaderKeyAccCtlReqMethod, "get")
	req.Header.Set(cors.HeaderKeyAccCtlReqHeaders, "Authorization, Content-Type")
	req.Header.Set("content-type", "application/json")

	w := httptest.NewRecorder()
	c.DoPreflight(w, req, func(response cors.PreflightResponse, w http.ResponseWriter, r *http.Request) {
		if response.HasError {
			t.Error("expected response.HasError to be false")
			return
		}
	})
}
