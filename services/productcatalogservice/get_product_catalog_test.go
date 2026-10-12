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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetProductReturnsCatalogEntry(t *testing.T) {
	catalog := &productCatalog{}
	if err := readCatalogFile(&catalog.catalog); err != nil {
		t.Fatal(err)
	}
	if len(catalog.catalog.Products) == 0 {
		t.Fatal("expected non-empty product catalog fixture")
	}
	wantID := catalog.catalog.Products[0].Id

	got, err := catalog.GetProduct(context.Background(), &pb.GetProductRequest{Id: wantID})
	if err != nil {
		t.Fatalf("GetProduct: %v", err)
	}
	if got.Id != wantID {
		t.Fatalf("expected id %q, got %q", wantID, got.Id)
	}
}

func TestGetProductMissingID(t *testing.T) {
	catalog := &productCatalog{}
	if err := readCatalogFile(&catalog.catalog); err != nil {
		t.Fatal(err)
	}

	_, err := catalog.GetProduct(context.Background(), &pb.GetProductRequest{Id: "does-not-exist"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", err)
	}
}
