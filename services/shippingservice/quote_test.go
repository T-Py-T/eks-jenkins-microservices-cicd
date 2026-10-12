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
	"regexp"
	"testing"
)

func TestCreateQuoteFromFloatSplitsUnitsAndCents(t *testing.T) {
	q := CreateQuoteFromFloat(8.99)
	if q.Dollars != 8 || q.Cents != 99 {
		t.Fatalf("CreateQuoteFromFloat(8.99) = %+v, want 8 dollars 99 cents", q)
	}
}

func TestCreateQuoteFromCountMatchesFlatRate(t *testing.T) {
	q := CreateQuoteFromCount(42)
	if q.Dollars != 8 || q.Cents != 99 {
		t.Fatalf("CreateQuoteFromCount should use flat 8.99 quote, got %+v", q)
	}
}

func TestCreateTrackingIdFormat(t *testing.T) {
	id := CreateTrackingId("Muffin Man, London, ")
	pattern := regexp.MustCompile(`^[A-Z]{2}-\d{5}-\d{9}$`)
	if !pattern.MatchString(id) {
		t.Fatalf("tracking id %q does not match expected pattern", id)
	}
	if len(id) != 18 {
		t.Fatalf("tracking id length = %d, want 18", len(id))
	}
}
