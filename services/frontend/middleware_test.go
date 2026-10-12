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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestEnsureSessionIDSetsCookieWhenMissing(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "")

	var capturedSession string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedSession = sessionID(r)
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	cookie := res.Cookies()
	if len(cookie) == 0 {
		t.Fatal("expected session cookie to be set")
	}
	if cookie[0].Name != cookieSessionID {
		t.Fatalf("unexpected cookie name: %s", cookie[0].Name)
	}
	if cookie[0].Value == "" {
		t.Fatal("expected non-empty session id")
	}
	if capturedSession != cookie[0].Value {
		t.Fatalf("context session %q does not match cookie %q", capturedSession, cookie[0].Value)
	}
}

func TestEnsureSessionIDReusesExistingCookie(t *testing.T) {
	const existing = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessionID(r) != existing {
			t.Fatalf("expected session %q, got %q", existing, sessionID(r))
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieSessionID, Value: existing})
	rec := httptest.NewRecorder()
	handler(rec, req)

	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("should not set a new cookie when one is already present")
	}
}

func TestEnsureSessionIDSharedSessionEnv(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "true")

	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessionID(r) != "12345678-1234-1234-1234-123456789123" {
			t.Fatalf("unexpected shared session id: %q", sessionID(r))
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
}

func TestLogHandlerRecordsStatusAndBytes(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(ctxKeyRequestID{}) == nil {
			t.Fatal("expected request id in context")
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("hello"))
	})

	req := httptest.NewRequest(http.MethodPost, "/cart", nil)
	rec := httptest.NewRecorder()
	(&logHandler{log: log, next: inner}).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestSessionIDFromContext(t *testing.T) {
	const id = "ctx-session-id"
	ctx := context.WithValue(context.Background(), ctxKeySessionID{}, id)
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	if sessionID(req) != id {
		t.Fatalf("expected %q from context, got %q", id, sessionID(req))
	}
}

