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

func TestValidPreflight(t *testing.T) {
	o := cors.Options{
		AllowedOrigins: []*cors.Match{
			cors.WC(`https?://.*\.theyakka\.com`), cors.WC(`https?://theyakka\.com`),
		},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		t.Error(err)
		return
	}

	req := buildPreflightRequest("https://theyakka.com")
	w := httptest.NewRecorder()
	c.ValidatePreflight(w, req, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
		if error != nil {
			t.Error("expected no errors to have occurred")
			return
		}
	})
}

func TestValidPreflightResponse(t *testing.T) {
	o := cors.Options{
		AllowedOrigins: []*cors.Match{
			cors.WC(`https?://.*\.theyakka\.com`), cors.WC(`https?://theyakka\.com`),
		},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		t.Error(err)
		return
	}

	req := buildPreflightRequest("https://theyakka.com")
	w := httptest.NewRecorder()
	c.ValidatePreflight(w, req, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
		if error != nil {
			t.Error("expected no errors to have occurred")
			return
		}
		header := w.Header()
		originOK := header.Get(cors.HeaderKeyAccCtlResAllowOrigin) == "https://theyakka.com"
		methodOK := header.Get(cors.HeaderKeyAccCtlResAllowMethods) == http.MethodGet
		headersOK := header.Get(cors.HeaderKeyAccCtlResAllowHeaders) == "Authorization, Content-Type"
		if !originOK || !methodOK || !headersOK {
			t.Error("response headers don't match expected")
		}
	})
}

func TestInvalidOriginPreflight(t *testing.T) {
	o := cors.Options{
		AllowedOrigins: []*cors.Match{
			cors.WC(`https?://.*\.theyakka\.com`), cors.WC(`https?://theyakka\.com`),
		},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		t.Error(err)
		return
	}

	req := buildPreflightRequest("https://google.com")
	w := httptest.NewRecorder()
	c.ValidatePreflight(w, req, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
		if error == nil {
			t.Error("expected an error to have occurred")
			return
		}
		if error.Code != cors.PreflightErrOriginNotAllowed {
			t.Error("expected error code to indicate the origin wasn't allowed")
			return
		}
	})
}

func TestValidOriginPreflightWithPort(t *testing.T) {
	o := cors.Options{
		AllowedOrigins: []*cors.Match{
			cors.WC(`https?://localhost:8[0-9]{3}`),
		},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		t.Error(err)
		return
	}

	req := buildPreflightRequest("http://localhost:8022")
	w := httptest.NewRecorder()
	c.ValidatePreflight(w, req, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
		if error != nil {
			t.Error("expected no error to have occurred")
			return
		}
	})
}

func BenchmarkExactPreflight(b *testing.B) {
	o := cors.Options{
		AllowedOrigins: []*cors.Match{
			cors.EM(`https://theyakka.com`),
		},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		b.Error(err)
		return
	}

	req := buildPreflightRequest("https://theyakka.com")
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.ValidatePreflight(w, req, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
		})
	}
}

func BenchmarkWildcardPreflight(b *testing.B) {
	o := cors.Options{
		AllowedOrigins: []*cors.Match{
			cors.WC(`https?://.*\.theyakka\.com`), cors.WC(`https?://theyakka\.com`),
		},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		b.Error(err)
		return
	}

	req := buildPreflightRequest("https://theyakka.com")
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.ValidatePreflight(w, req, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
		})
	}
}

func BenchmarkWildcardPortPreflight(b *testing.B) {
	o := cors.Options{
		AllowedOrigins: []*cors.Match{
			cors.WC(`https?://localhost:8[0-9]{3}`),
		},
		AllowedHeaders: cors.DefaultHeadersWith("Authorization"),
	}
	c, err := o.NewCORS()
	if err != nil {
		b.Error(err)
		return
	}

	req := buildPreflightRequest("http://localhost:8022")
	w := httptest.NewRecorder()
	handler := fakeHandler(c)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
	}
}

func TestAllowAllPreflight(t *testing.T) {
	c, err := cors.AllowAll()
	if err != nil {
		t.Error(err)
		return
	}

	req := buildPreflightRequest("https://somethingdodgy.xxx")
	w := httptest.NewRecorder()
	c.ValidatePreflight(w, req, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
		if error != nil {
			t.Error("expected no error to have occurred")
			return
		}
	})
}

func buildPreflightRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodOptions, url, nil)
	req.Header.Set(cors.HeaderKeyReqOrigin, url)
	req.Header.Set(cors.HeaderKeyAccCtlReqMethod, "get")
	req.Header.Set(cors.HeaderKeyAccCtlReqHeaders, "Authorization, Content-Type")
	req.Header.Set("content-type", "application/json")
	return req
}

func fakeHandler(c *cors.CORS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.ValidatePreflight(w, r, func(w http.ResponseWriter, r *http.Request, error *cors.ValidationError) {
			_, _ = w.Write([]byte("ok"))
		})
	})
}
