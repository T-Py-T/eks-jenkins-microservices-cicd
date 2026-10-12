// Copyright 2023 Google LLC
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

func TestSearchProductsCaseInsensitive(t *testing.T) {
	catalog := &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "1", Name: "Sunglasses", Description: "Aviator style"},
				{Id: "2", Name: "Watch", Description: "Gold stainless steel"},
			},
		},
	}

	resp, err := catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: "GOLD"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 1 || resp.Results[0].Id != "2" {
		t.Fatalf("got %#v, want watch product", resp.Results)
	}
}

func TestSearchProductsMatchesDescription(t *testing.T) {
	catalog := &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "1", Name: "Mug", Description: "Ceramic kitchen mug"},
			},
		},
	}

	resp, err := catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: "kitchen"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected one match, got %d", len(resp.Results))
	}
}

func TestSearchProductsNoMatch(t *testing.T) {
	catalog := &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "1", Name: "Mug", Description: "Ceramic kitchen mug"},
			},
		},
	}

	resp, err := catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: "loafers"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 0 {
		t.Fatalf("expected no matches, got %d", len(resp.Results))
	}
}
