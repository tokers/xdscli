// Copyright 2020 xdscli Authors.
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

package main

import (
	"fmt"
	"runtime"
)

var (
	_gitSHA    = "unknown"
	_version   = "0.0.1"
	_goVersion = runtime.Version()
	_goOSArch  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

func showVersionAndQuit() {
	fmt.Printf("xdscli Version: %s\n", _version)
	fmt.Printf("Git SHA: %s\n", _gitSHA)
	fmt.Printf("Go Version: %s\n", _goVersion)
	fmt.Printf("Go OS/Arch: %s\n", _goOSArch)
	exitWithError(_exitSuccess, nil)
}
