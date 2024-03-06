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

// Basic main function to test hosttracker.
// Example usage: go run main/main.go google.com

package main

import (
	"log"
	"time"
	"os"

	"github.com/bjornleffler/k8s-golang-utils/servicetracker"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s HOSTNAME", os.Args[0])
		os.Exit(1)
	}
	hostname := os.Args[1]
	_ = servicetracker.MakeTracker(hostname)
	for {
		time.Sleep(1 * time.Second)
	}
}
