// Copyright 2016 Michael Stapelberg and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package dispatch offers a registry in which scan actions are resolved. E.g.,
// as soon as an AirScan-compatible scanner is discovered, the default scan
// action is hooked up to it.
package dispatch

import (
	"fmt"
	"sync"
)

type ScanAction interface {
	TargetName() string
	Scan(user string) (name string, _ error)
}

var (
	actionsMu sync.Mutex
	actions   []ScanAction
)

func Register(a ScanAction) {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	actions = append(actions, a)
}

func Scan(user string) (name string, _ error) {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	if len(actions) == 0 {
		// should be reported as sth along the lines of “no scanners found”
		return "", fmt.Errorf("dispatch failed: no actions registered")
	}
	mostRecent := actions[len(actions)-1]
	return mostRecent.Scan(user)
}
