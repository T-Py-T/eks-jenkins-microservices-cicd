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
	"errors"
	"net"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
)

type contextAwareCurrencyService struct{}

func (s *contextAwareCurrencyService) GetSupportedCurrencies(context.Context, *pb.Empty) (*pb.GetSupportedCurrenciesResponse, error) {
	return &pb.GetSupportedCurrenciesResponse{}, nil
}

func (s *contextAwareCurrencyService) Convert(ctx context.Context, req *pb.CurrencyConversionRequest) (*pb.Money, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &pb.Money{
		CurrencyCode: req.ToCode,
		Units:        req.From.GetUnits(),
		Nanos:        req.From.GetNanos(),
	}, nil
}

func TestConvertCurrencyPropagatesCancelledContext(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCurrencyServiceServer(grpcServer, &contextAwareCurrencyService{})
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc client: %v", err)
	}
	defer conn.Close()

	cs := &checkoutService{currencySvcConn: conn}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = cs.convertCurrency(ctx, &pb.Money{CurrencyCode: "USD", Units: 10}, "EUR")
	if err == nil {
		t.Fatal("expected error when request context is already cancelled")
	}
	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("expected context cancellation error, got: %v", err)
	}
}
