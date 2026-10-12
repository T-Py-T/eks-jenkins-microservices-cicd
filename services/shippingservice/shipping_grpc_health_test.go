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
	"testing"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestShippingServiceHealthCheckServing(t *testing.T) {
	s := &server{}
	resp, err := s.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if resp.Status != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("status %v, want SERVING", resp.Status)
	}
}

func TestShippingServiceHealthWatchUnimplemented(t *testing.T) {
	s := &server{}
	err := s.Watch(&healthpb.HealthCheckRequest{}, nil)
	if err == nil {
		t.Fatal("expected Watch to return an error")
	}
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("got code %v, want Unimplemented", status.Code(err))
	}
}
