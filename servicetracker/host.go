// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package servicetracker

import (
	"net/netip"
	"time"
)

type HostStatus int

const (
	Active HostStatus = iota
	Removed
)

var (
	Timeout = time.Duration(60 * time.Second)
)

type Host struct {
	Addr        netip.Addr
	Status      HostStatus
	LastActive  time.Time
}

func MakeHost(addr netip.Addr) *Host {
	h := &Host{
		Addr:        addr,
		Status:      Active,
		LastActive:  time.Now(),
	}
	return h
}

func (h *Host) IsActive() bool {
	return h.Status == Active
}

func (h *Host) IsRemoved() bool {
	return h.Status == Removed
}

// SetActive sets the Host status to active. Returns true on status changes.
func (h *Host) SetActive() bool {
	updated := (h.Status != Active)
	h.Status = Active
	h.LastActive = time.Now()
	return updated
}

func (h *Host) SetRemoved() {
	h.Status = Removed
}
