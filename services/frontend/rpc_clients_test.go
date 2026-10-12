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
	"context"
	"net"
	"strings"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufConnSize = 1024 * 1024

func dialBufConn(lis *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func startBufGRPC(register func(*grpc.Server)) (*grpc.ClientConn, func()) {
	lis := bufconn.Listen(bufConnSize)
	s := grpc.NewServer()
	register(s)
	go func() { _ = s.Serve(lis) }()

	conn, err := dialBufConn(lis)
	if err != nil {
		s.Stop()
		panic(err)
	}
	cleanup := func() {
		_ = conn.Close()
		s.Stop()
	}
	return conn, cleanup
}

type stubCurrencyService struct {
	codes   []string
	convert func(context.Context, *pb.CurrencyConversionRequest) (*pb.Money, error)
}

func (s *stubCurrencyService) GetSupportedCurrencies(context.Context, *pb.Empty) (*pb.GetSupportedCurrenciesResponse, error) {
	return &pb.GetSupportedCurrenciesResponse{CurrencyCodes: s.codes}, nil
}

func (s *stubCurrencyService) Convert(ctx context.Context, req *pb.CurrencyConversionRequest) (*pb.Money, error) {
	if s.convert != nil {
		return s.convert(ctx, req)
	}
	return req.From, nil
}

type stubProductCatalog struct {
	products map[string]*pb.Product
}

func (s *stubProductCatalog) ListProducts(context.Context, *pb.Empty) (*pb.ListProductsResponse, error) {
	return &pb.ListProductsResponse{}, nil
}

func (s *stubProductCatalog) GetProduct(_ context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	p, ok := s.products[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "missing product %s", req.Id)
	}
	return p, nil
}

func (s *stubProductCatalog) SearchProducts(context.Context, *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	return &pb.SearchProductsResponse{}, nil
}

type stubRecommendationService struct {
	productIDs []string
}

func (s *stubRecommendationService) ListRecommendations(context.Context, *pb.ListRecommendationsRequest) (*pb.ListRecommendationsResponse, error) {
	return &pb.ListRecommendationsResponse{ProductIds: s.productIDs}, nil
}

type stubShippingService struct {
	quoteUSD *pb.Money
}

func (s *stubShippingService) GetQuote(context.Context, *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	return &pb.GetQuoteResponse{CostUsd: s.quoteUSD}, nil
}

func (s *stubShippingService) ShipOrder(context.Context, *pb.ShipOrderRequest) (*pb.ShipOrderResponse, error) {
	return &pb.ShipOrderResponse{}, nil
}

func TestGetCurrenciesFiltersToWhitelist(t *testing.T) {
	conn, cleanup := startBufGRPC(func(s *grpc.Server) {
		pb.RegisterCurrencyServiceServer(s, &stubCurrencyService{
			codes: []string{"USD", "XXX", "EUR", "BTC", "GBP"},
		})
	})
	defer cleanup()

	fe := &frontendServer{currencySvcConn: conn}
	got, err := fe.getCurrencies(context.Background())
	if err != nil {
		t.Fatalf("getCurrencies: %v", err)
	}

	want := map[string]bool{"USD": true, "EUR": true, "GBP": true}
	if len(got) != len(want) {
		t.Fatalf("got %d currencies %v, want %d whitelisted", len(got), got, len(want))
	}
	for _, code := range got {
		if !want[code] {
			t.Fatalf("unexpected currency %q in %v", code, got)
		}
	}
}

func TestGetRecommendationsTruncatesToFourProducts(t *testing.T) {
	ids := []string{"p1", "p2", "p3", "p4", "p5", "p6"}
	products := make(map[string]*pb.Product, len(ids))
	for _, id := range ids {
		products[id] = &pb.Product{Id: id, Name: id}
	}

	catalogConn, catalogCleanup := startBufGRPC(func(s *grpc.Server) {
		pb.RegisterProductCatalogServiceServer(s, &stubProductCatalog{products: products})
	})
	defer catalogCleanup()

	recConn, recCleanup := startBufGRPC(func(s *grpc.Server) {
		pb.RegisterRecommendationServiceServer(s, &stubRecommendationService{productIDs: ids})
	})
	defer recCleanup()

	fe := &frontendServer{
		productCatalogSvcConn: catalogConn,
		recommendationSvcConn: recConn,
	}

	got, err := fe.getRecommendations(context.Background(), "user-1", []string{"in-cart"})
	if err != nil {
		t.Fatalf("getRecommendations: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("got %d products, want UI cap of 4", len(got))
	}
	for i, p := range got {
		if p.Id != ids[i] {
			t.Fatalf("position %d: got id %q, want %q", i, p.Id, ids[i])
		}
	}
}

func TestGetShippingQuoteWrapsCurrencyConversionFailure(t *testing.T) {
	shipConn, shipCleanup := startBufGRPC(func(s *grpc.Server) {
		pb.RegisterShippingServiceServer(s, &stubShippingService{
			quoteUSD: &pb.Money{CurrencyCode: "USD", Units: 8, Nanos: 990000000},
		})
	})
	defer shipCleanup()

	currConn, currCleanup := startBufGRPC(func(s *grpc.Server) {
		pb.RegisterCurrencyServiceServer(s, &stubCurrencyService{
			convert: func(context.Context, *pb.CurrencyConversionRequest) (*pb.Money, error) {
				return nil, status.Error(codes.Unavailable, "conversion unavailable")
			},
		})
	})
	defer currCleanup()

	fe := &frontendServer{
		shippingSvcConn: shipConn,
		currencySvcConn: currConn,
	}

	_, err := fe.getShippingQuote(context.Background(), []*pb.CartItem{{ProductId: "1", Quantity: 1}}, "EUR")
	if err == nil {
		t.Fatal("expected error when currency conversion fails")
	}
	if !strings.Contains(err.Error(), "failed to convert currency for shipping cost") {
		t.Fatalf("expected wrapped error, got %q", err.Error())
	}
}
