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

func TestCreateQuoteFromFloat_avoidsBinaryFractionTruncation(t *testing.T) {
	tests := []struct {
		value         float64
		wantDollars   uint32
		wantCents     uint32
	}{
		{2.51, 2, 51},
		{19.99, 19, 99},
		{8.99, 8, 99},
	}
	for _, tt := range tests {
		quote := CreateQuoteFromFloat(tt.value)
		if quote.Dollars != tt.wantDollars || quote.Cents != tt.wantCents {
			t.Errorf(
				"CreateQuoteFromFloat(%v) = $%d.%02d, want $%d.%02d",
				tt.value,
				quote.Dollars,
				quote.Cents,
				tt.wantDollars,
				tt.wantCents,
			)
		}
	}
}
