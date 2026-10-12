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

import "testing"

func TestSetPlatformDetails(t *testing.T) {
	cases := []struct {
		env      string
		provider string
		css      string
	}{
		{"aws", "AWS", "aws-platform"},
		{"onprem", "On-Premises", "onprem-platform"},
		{"azure", "Azure", "azure-platform"},
		{"gcp", "Google Cloud", "gcp-platform"},
		{"alibaba", "Alibaba Cloud", "alibaba-platform"},
		{"", "local", "local"},
		{"unknown-cloud", "local", "local"},
	}

	for _, tc := range cases {
		plat := &platformDetails{}
		plat.setPlatformDetails(tc.env)
		if plat.provider != tc.provider || plat.css != tc.css {
			t.Fatalf("env %q: got provider=%q css=%q, want provider=%q css=%q",
				tc.env, plat.provider, plat.css, tc.provider, tc.css)
		}
	}
}
