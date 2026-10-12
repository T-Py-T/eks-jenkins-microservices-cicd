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

func TestCreateQuoteFromFloat(t *testing.T) {
	quote := CreateQuoteFromFloat(8.99)
	if quote.Dollars != 8 || quote.Cents != 99 {
		t.Fatalf("got $%d.%02d, want $8.99", quote.Dollars, quote.Cents)
	}
	if quote.String() != "$8.99" {
		t.Fatalf("String() = %q, want $8.99", quote.String())
	}
}

func TestCreateQuoteFromCountUsesFlatRate(t *testing.T) {
	for _, count := range []int{0, 1, 99} {
		quote := CreateQuoteFromCount(count)
		if quote.Dollars != 8 || quote.Cents != 99 {
			t.Fatalf("count %d: got $%d.%02d, want flat $8.99", count, quote.Dollars, quote.Cents)
		}
	}
}

func TestCreateQuoteFromFloatZero(t *testing.T) {
	quote := CreateQuoteFromFloat(0)
	if quote.Dollars != 0 || quote.Cents != 0 {
		t.Fatalf("got $%d.%02d, want $0.00", quote.Dollars, quote.Cents)
	}
}
