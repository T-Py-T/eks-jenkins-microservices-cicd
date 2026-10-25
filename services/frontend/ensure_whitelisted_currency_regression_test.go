// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnsureWhitelistedCurrencyResetsUnsupportedCode(t *testing.T) {
	var served bool
	handler := ensureWhitelistedCurrency(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		served = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieCurrency, Value: "BOGUS"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !served {
		t.Fatal("expected inner handler to run")
	}

	resp := rec.Result()
	var reset *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == cookieCurrency {
			reset = c
			break
		}
	}
	if reset == nil {
		t.Fatal("expected Set-Cookie for shop_currency")
	}
	if reset.Value != defaultCurrency {
		t.Fatalf("expected currency reset to %q, got %q", defaultCurrency, reset.Value)
	}
}

func TestEnsureWhitelistedCurrencyKeepsSupportedCode(t *testing.T) {
	handler := ensureWhitelistedCurrency(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieCurrency, Value: "EUR"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	for _, c := range rec.Result().Cookies() {
		if c.Name == cookieCurrency {
			t.Fatalf("did not expect currency cookie to be rewritten, got %q", c.Value)
		}
	}
}
