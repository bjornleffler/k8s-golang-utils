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

// Utility for tracking a list of hosts resolved from a service hostname.
// Intended to track IPs in a K8s daemonset, etc.

package servicetracker

import (
	"log"
	"net"
	"net/netip"
	"sort"
	"time"
)

const (
	SecondsBetweenLookups = 10
)

type Tracker struct {
	hostname   string
	hosts      map[netip.Addr]*Host
	LastUpdate time.Time
}

func MakeTracker(hostname string) *Tracker {
	t := &Tracker{
		hostname:   hostname,
		hosts:      make(map[netip.Addr]*Host),
		LastUpdate: time.Now(),
	}
	go t.track()
	return t
}

func (t *Tracker) track() {
	for {
		if ips, err := net.LookupIP(t.hostname); err == nil {
			// 1. Update active hosts.
			active := map[netip.Addr]bool{}
			for _, ip := range ips {
				addr, ok := netip.AddrFromSlice(ip)
				if !ok {
					log.Printf("Failed to convert IP: %v", ip)
					continue
				}
				active[addr] = true
				if host, ok := t.hosts[addr]; !ok {
					t.AddHost(addr)
				} else {
					t.UpdateHost(host)
				}
			}

			// 2. Mark other hosts as removed.
			for _, host := range t.hosts {
				if _, ok := active[host.Addr]; ok {
					continue
				}
				if !host.IsRemoved() {
					t.RemoveHost(host)
				}
			}
		}
		time.Sleep(SecondsBetweenLookups * time.Second)
	}
}

func (t *Tracker) AddHost(addr netip.Addr) {
	// log.Printf("Tracker: new host: %v", addr)
	t.hosts[addr] = MakeHost(addr)
	t.LastUpdate = time.Now()
}

func (t *Tracker) UpdateHost(host *Host) {
	if updated := host.SetActive(); updated {
		t.LastUpdate = time.Now()
	}
}

func (t *Tracker) RemoveHost(host *Host) {
	// log.Printf("Marking host %v as removed.", host.Addr)
	host.SetRemoved()
}

// ActiveHosts returns the list of active hosts in sorted order.
func (t *Tracker) ActiveHosts() []*Host {
	result := []*Host{}
	for _, host := range t.hosts {
		if host.IsActive() {
			result = append(result, host)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Addr.Less(result[j].Addr)
	})
	return result
}

func (t *Tracker) RemovedHosts() []*Host {
	result := []*Host{}
	for _, host := range t.hosts {
		if host.IsRemoved() {
			result = append(result, host)
		}
	}
	return result
}
