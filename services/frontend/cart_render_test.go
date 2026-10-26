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

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func TestCartSizeSumsQuantities(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "a", Quantity: 2},
		{ProductId: "b", Quantity: 3},
	}
	if cartSize(items) != 5 {
		t.Fatalf("expected cart size 5, got %d", cartSize(items))
	}
	if cartSize(nil) != 0 {
		t.Fatalf("expected empty cart size 0, got %d", cartSize(nil))
	}
}

func TestCartIDsPreservesOrder(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "first"},
		{ProductId: "second"},
	}
	got := cartIDs(items)
	want := []string{"first", "second"}
	if len(got) != len(want) {
		t.Fatalf("expected %d ids, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestRenderMoneyAndCurrencyLogos(t *testing.T) {
	m := pb.Money{CurrencyCode: "EUR", Units: 12, Nanos: 250_000_000}
	if renderMoney(m) != "€12.25" {
		t.Fatalf("unexpected renderMoney output: %q", renderMoney(m))
	}
	if renderCurrencyLogo("JPY") != "¥" {
		t.Fatalf("expected yen logo")
	}
	if renderCurrencyLogo("UNKNOWN") != "$" {
		t.Fatalf("expected default dollar logo for unknown currency")
	}
}

func TestStringInSlice(t *testing.T) {
	slice := []string{"USD", "EUR", "GBP"}
	if !stringinSlice(slice, "EUR") {
		t.Fatal("expected EUR to be found")
	}
	if stringinSlice(slice, "CAD") {
		t.Fatal("expected CAD to be absent")
	}
}

func TestCurrentCurrencyCookieAndDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if currentCurrency(req) != defaultCurrency {
		t.Fatalf("expected default %q, got %q", defaultCurrency, currentCurrency(req))
	}

	req.AddCookie(&http.Cookie{Name: cookieCurrency, Value: "GBP"})
	if currentCurrency(req) != "GBP" {
		t.Fatalf("expected GBP from cookie, got %q", currentCurrency(req))
	}
}
