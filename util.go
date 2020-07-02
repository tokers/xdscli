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
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	_struct "github.com/golang/protobuf/ptypes/struct"
)

const (
	_apiVersion2          = "v2"
	_serviceNodeSeparator = "~"
)

var (
	_typeURLMap = map[string]string{
		"eds": "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
		"cds": "type.googleapis.com/envoy.api.v2.Cluster",
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func validateAPIVersion(ver string) error {
	switch ver {
	case _apiVersion2:
	default:
		return fmt.Errorf("bad api version: %s", ver)
	}
	return nil
}

func validateTimeoutValue() error {
	if _gFlags.dialTimeout < 0 {
		return _errInvalidDialTimeout
	}
	if _gFlags.readTimeout < 0 {
		return _errInvalidReadTimeout
	}
	if _gFlags.sendTimeout < 0 {
		return _errInvalidSendTimeout
	}
	return nil
}

func validateOutputFormat() error {
	format := strings.ToLower(_gFlags.outputFormat)
	switch format {
	case "json", "yaml", "simple":
		_gFlags.outputFormat = strings.ToLower(format)
	default:
		return _errInvalidOutputFormat
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

func buildOutputMarshaller(format string) marshaller {
	switch format {
	case "json":
		return newJSONMarshaller()
	case "simple":
		return newDefaultMarshaller()
	case "yaml":
		return newYAMLMarshaller()
	default:
		panic("not implemented yet")
	}
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

func validateXDS() error {
	// FIXME Maybe just let user ensure the validity of node id?
	if _gFlags.xds.node != "" {
		parts := strings.Split(_gFlags.xds.node, _serviceNodeSeparator)
		if len(parts) != 4 {
			return _errInvalidNode
		}
		if parts[0] != "sidecar" && parts[0] != "router" {
			return _errInvalidNode
		}
	} else {
		_gFlags.xds.node = genNodeID()
	}

	return nil
}

func buildNodeMetadata(meta string) (*_struct.Struct, error) {
	metadata := &_struct.Struct{
		Fields: make(map[string]*_struct.Value),
	}
	for _, part := range strings.Split(meta, ",") {
		subpart := strings.Split(part, "=")
		if len(subpart) != 2 {
			return nil, _errInvalidNodeMetaFormat
		}
		metadata.Fields[subpart[0]] = &_struct.Value{
			Kind: &_struct.Value_StringValue{
				StringValue: subpart[1],
			},
		}
	}
	return metadata, nil
}

func genNodeID() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	nics, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	addr := "0.0.0.0"
	for _, nic := range nics {
		// Find the first non loopback NIC.
		if nic.Name == "lo" {
			continue
		}

		addrs, err := nic.Addrs()
		if err != nil {
			panic(err)
		}

		if len(addrs) > 0 {
			addr = strings.Split(addrs[0].String(), "/")[0]
			break
		}
	}

	return fmt.Sprintf("sidecar~%s~%d~%s", addr, rand.Int(), hostname)
}

func validateOptions() error {
	if err := validateAPIVersion(_gFlags.xds.apiVersion); err != nil {
		return err
	}

	if err := validateTimeoutValue(); err != nil {
		return err
	}

	if _gFlags.grpcMaxCallRecvSize <= 0 {
		return _errInvalidGRPCMaxCallRecvSize
	}

	if err := validateOutputFormat(); err != nil {
		return err
	}

	if err := validateXDS(); err != nil {
		return err
	}
	return nil
}
