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
	"testing"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

type emptyAdService struct{}

func (s *emptyAdService) GetAds(context.Context, *pb.AdRequest) (*pb.AdResponse, error) {
	return &pb.AdResponse{}, nil
}

func TestChooseAdEmptyResponseDoesNotPanic(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAdServiceServer(grpcServer, &emptyAdService{})
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc client: %v", err)
	}
	defer conn.Close()

	fe := &frontendServer{adSvcConn: conn}
	log := logrus.New()

	var ad *pb.Ad
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("chooseAd panicked on empty ad list: %v", r)
			}
		}()
		ad = fe.chooseAd(context.Background(), nil, log)
	}()

	if ad != nil {
		t.Fatalf("expected nil ad when ad service returns no ads, got %+v", ad)
	}
}
