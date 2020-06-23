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
	"net"
)

const (
	_apiVersion2 = "v2"
)

var (
	_typeURLMap = map[string]string{
		"eds": "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
	}
)

func validateAPIVersion(ver string) error {
	switch ver {
	case _apiVersion2:
	default:
		return fmt.Errorf("bad api version: %s", ver)
	}
	return nil
}

func validateAndResolveServers(servers []string) ([]string, error) {
	var endpoints []string
	for _, srv := range servers {
		host, port, err := net.SplitHostPort(srv)
		if err != nil {
			return nil, err
		}

		if ip := net.ParseIP(host); ip != nil {
			endpoints = append(endpoints, net.JoinHostPort(ip.String(), port))
		} else {
			// Try to resolve this host.
			// TODO Support custom DNS resolver by adding a new command option
			// like --resolver.
			addrs, err := net.LookupHost(host)
			if err != nil {
				return nil, err
			}

			for _, addr := range addrs {
				endpoints = append(endpoints, net.JoinHostPort(addr, port))
			}
		}
	}

	return endpoints, nil
}

func getDiscoveryServiceTypeUrl(apiVersion, ds string) (string, error) {
	// TODO support api version 3.
	if apiVersion != _apiVersion2 {
		panic("unknown api version")
	}
	typeURL, ok := _typeURLMap[ds]
	if !ok {
		return "", fmt.Errorf("unknown discovery service: %s", ds)
	}
	return typeURL, nil
}
