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
	"net"
	"sync"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type stubRecommendationService struct {
	pb.RecommendationServiceServer
	productIDs []string
}

func (s *stubRecommendationService) ListRecommendations(context.Context, *pb.ListRecommendationsRequest) (*pb.ListRecommendationsResponse, error) {
	return &pb.ListRecommendationsResponse{ProductIds: s.productIDs}, nil
}

type stubProductCatalogService struct {
	pb.ProductCatalogServiceServer
	mu    sync.Mutex
	calls int
}

func (s *stubProductCatalogService) GetProduct(context.Context, *pb.GetProductRequest) (*pb.Product, error) {
	s.mu.Lock()
	s.calls++
	s.mu.Unlock()
	return &pb.Product{Id: "stub"}, nil
}

func (s *stubProductCatalogService) getProductCalls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func startTestGRPCServer(t *testing.T, register func(*grpc.Server)) (*grpc.ClientConn, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	server := grpc.NewServer()
	register(server)

	go server.Serve(listener)

	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	cleanup := func() {
		conn.Close()
		server.Stop()
	}
	return conn, cleanup
}

func TestGetRecommendationsFetchesAtMostFourProducts(t *testing.T) {
	recommendedIDs := []string{"p1", "p2", "p3", "p4", "p5", "p6"}
	recConn, stopRec := startTestGRPCServer(t, func(s *grpc.Server) {
		pb.RegisterRecommendationServiceServer(s, &stubRecommendationService{productIDs: recommendedIDs})
	})
	defer stopRec()

	catalog := &stubProductCatalogService{}
	catalogConn, stopCatalog := startTestGRPCServer(t, func(s *grpc.Server) {
		pb.RegisterProductCatalogServiceServer(s, catalog)
	})
	defer stopCatalog()

	fe := &frontendServer{
		recommendationSvcConn: recConn,
		productCatalogSvcConn: catalogConn,
	}

	products, err := fe.getRecommendations(context.Background(), "user", nil)
	if err != nil {
		t.Fatalf("getRecommendations: %v", err)
	}
	if len(products) != 4 {
		t.Fatalf("expected 4 products, got %d", len(products))
	}
	if catalog.getProductCalls() != 4 {
		t.Fatalf("expected 4 catalog lookups, got %d", catalog.getProductCalls())
	}
}
