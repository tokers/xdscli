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
	"time"
)

// xdsFlags are flags the defined about the
// DiscoveryRequest/DeltaDiscoveryRequest.
type xdsFlags struct {
	node               string
	nodeMetadata       string
	initialVersionInfo string
	errorDetail        string
	resourceNames      []string
	apiVersion         string
}

// globalFlags are flags that defined globally and are inherited to all
// sub-commands.
type globalFlags struct {
	xds xdsFlags

	dialTimeout time.Duration
	readTimeout time.Duration
	sendTimeout time.Duration
	timeout     time.Duration

	outputFormat string
	servers      []string
	wait         bool
	showVersion  bool
}
