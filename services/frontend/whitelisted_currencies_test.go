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

import "testing"

func filterWhitelistedCurrencyCodes(codes []string) []string {
	var out []string
	for _, code := range codes {
		if whitelistedCurrencies[code] {
			out = append(out, code)
		}
	}
	return out
}

func TestWhitelistedCurrenciesContainsCheckoutDefaults(t *testing.T) {
	required := []string{"USD", "EUR", "CAD", "JPY", "GBP", "TRY"}
	for _, code := range required {
		if !whitelistedCurrencies[code] {
			t.Fatalf("expected %q to be whitelisted", code)
		}
	}
}

func TestFilterWhitelistedCurrencyCodes(t *testing.T) {
	input := []string{"USD", "BTC", "EUR", "CHF", "JPY"}
	got := filterWhitelistedCurrencyCodes(input)
	want := []string{"USD", "EUR", "JPY"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestFilterWhitelistedCurrencyCodesEmptyInput(t *testing.T) {
	if len(filterWhitelistedCurrencyCodes(nil)) != 0 {
		t.Fatal("expected empty result for nil input")
	}
}
