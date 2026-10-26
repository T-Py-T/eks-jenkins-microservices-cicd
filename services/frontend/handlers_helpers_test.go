// Copyright 2018 Google LLC
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
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func TestCartSizeSumsItemQuantities(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "a", Quantity: 2},
		{ProductId: "b", Quantity: 3},
	}
	if got := cartSize(items); got != 5 {
		t.Fatalf("cartSize() = %d, want 5", got)
	}
}

func TestCartIDsReturnsProductIdsInOrder(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "first"},
		{ProductId: "second"},
	}
	got := cartIDs(items)
	want := []string{"first", "second"}
	if len(got) != len(want) {
		t.Fatalf("cartIDs() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("cartIDs()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestStringInSlice(t *testing.T) {
	envs := []string{"local", "gcp", "aws"}
	if !stringinSlice(envs, "gcp") {
		t.Fatal("expected gcp to be in slice")
	}
	if stringinSlice(envs, "invalid") {
		t.Fatal("expected invalid to be absent from slice")
	}
}

func TestRenderMoneyAndCurrencyLogo(t *testing.T) {
	m := pb.Money{CurrencyCode: "EUR", Units: 12, Nanos: 340000000}
	rendered := renderMoney(m)
	if !strings.HasPrefix(rendered, "€") {
		t.Fatalf("renderMoney() = %q, want euro prefix", rendered)
	}
	if !strings.Contains(rendered, "12.34") {
		t.Fatalf("renderMoney() = %q, want dollars.cents formatting", rendered)
	}
	if renderCurrencyLogo("JPY") != "¥" {
		t.Fatalf("unexpected JPY logo")
	}
	if renderCurrencyLogo("UNKNOWN") != "$" {
		t.Fatalf("expected default logo for unknown currency")
	}
}

func TestCurrentCurrencyUsesCookieOrDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if currentCurrency(req) != defaultCurrency {
		t.Fatalf("default currency = %q, want %q", currentCurrency(req), defaultCurrency)
	}

	req.AddCookie(&http.Cookie{Name: cookieCurrency, Value: "GBP"})
	if currentCurrency(req) != "GBP" {
		t.Fatalf("cookie currency = %q, want GBP", currentCurrency(req))
	}
}

func TestSetCurrencyHandlerSetsCookieAndRedirects(t *testing.T) {
	fe := &frontendServer{}
	log := logrus.New()
	req := httptest.NewRequest(http.MethodPost, "/setCurrency", strings.NewReader("currency_code=EUR"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := req.ParseForm(); err != nil {
		t.Fatalf("ParseForm: %v", err)
	}
	req = req.WithContext(context.WithValue(req.Context(), ctxKeyLog{}, log))
	rr := httptest.NewRecorder()

	fe.setCurrencyHandler(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusFound)
	}
	if loc := rr.Header().Get("Location"); loc != "/" {
		t.Fatalf("Location = %q, want /", loc)
	}
	res := rr.Result()
	defer res.Body.Close()
	cookies := res.Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == cookieCurrency && c.Value == "EUR" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected currency cookie to be set to EUR")
	}
}

func TestSessionIDFromContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if sessionID(req) != "" {
		t.Fatal("expected empty session when context unset")
	}
	ctx := context.WithValue(req.Context(), ctxKeySessionID{}, "sess-abc")
	req = req.WithContext(ctx)
	if sessionID(req) != "sess-abc" {
		t.Fatalf("sessionID() = %q, want sess-abc", sessionID(req))
	}
}
