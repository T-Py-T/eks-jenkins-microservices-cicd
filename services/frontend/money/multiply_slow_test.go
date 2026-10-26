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

package money

import "testing"

func TestMultiplySlow_zeroQuantity(t *testing.T) {
	got := MultiplySlow(mmc(12, 500000000, "USD"), 0)
	want := mmc(0, 0, "USD")
	if !AreEquals(got, want) {
		t.Fatalf("MultiplySlow(..., 0) = %v, want %v", got, want)
	}
}
