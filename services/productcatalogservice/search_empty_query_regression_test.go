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
	"context"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
)

func TestSearchProductsEmptyQueryReturnsNoResults(t *testing.T) {
	resp, err := mockProductCatalog.SearchProducts(context.Background(),
		&pb.SearchProductsRequest{Query: ""},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 0 {
		t.Fatalf("empty search query should not match every product, got %d results", len(resp.Results))
	}
}

func TestSearchProductsWhitespaceQueryReturnsNoResults(t *testing.T) {
	resp, err := mockProductCatalog.SearchProducts(context.Background(),
		&pb.SearchProductsRequest{Query: "   \t"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 0 {
		t.Fatalf("whitespace-only search query should return no results, got %d", len(resp.Results))
	}
}
