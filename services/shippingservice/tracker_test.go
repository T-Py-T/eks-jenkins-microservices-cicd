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
	"regexp"
	"strings"
	"testing"
)

var trackingIDPattern = regexp.MustCompile(`^[A-Z]{2}-\d+-\d+$`)

func TestCreateTrackingIdFormat(t *testing.T) {
	salt := "Muffin Man, London, "
	for i := 0; i < 25; i++ {
		id := CreateTrackingId(salt)
		if !trackingIDPattern.MatchString(id) {
			t.Fatalf("tracking id %q does not match expected pattern", id)
		}
	}
}

func TestCreateTrackingIdReflectsSaltLength(t *testing.T) {
	saltA := "Muffin Man, London, "
	saltB := "abcdefghij"

	idA := CreateTrackingId(saltA)
	idB := CreateTrackingId(saltB)

	if trackingSaltLength(idA) != len(saltA) {
		t.Fatalf("id %q should encode salt length %d", idA, len(saltA))
	}
	if trackingSaltLength(idB) != len(saltB) {
		t.Fatalf("id %q should encode salt length %d", idB, len(saltB))
	}
}

func trackingSaltLength(id string) int {
	// Format: LL-<len><3digits>-<len/2><7digits>
	parts := strings.Split(id, "-")
	if len(parts) != 3 || len(parts[1]) < 4 {
		return -1
	}
	encodedLen := parts[1][:len(parts[1])-3]
	value := 0
	for _, ch := range encodedLen {
		if ch < '0' || ch > '9' {
			return -1
		}
		value = value*10 + int(ch-'0')
	}
	return value
}
