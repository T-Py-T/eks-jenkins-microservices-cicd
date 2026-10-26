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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestCheckoutServiceHealthCheck(t *testing.T) {
	svc := &checkoutService{}
	resp, err := svc.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if resp.Status != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("status %v, want SERVING", resp.Status)
	}
}

func TestCheckoutServiceHealthWatchUnimplemented(t *testing.T) {
	svc := &checkoutService{}
	err := svc.Watch(&healthpb.HealthCheckRequest{}, nil)
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("got code %v, want Unimplemented", status.Code(err))
	}
}
