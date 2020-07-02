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
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"

	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gogo/protobuf/proto"
)

type discoveryResponse struct {
	VersionInfo  string             `json:"version_info,omitempty",yaml:"version_info,omitempty"`
	Resources    []interface{}      `json:"resources,omitempty",yaml:"resources,omitempty"`
	Canary       bool               `json:"canary,omitempty",yaml:"canary,omitempty"`
	TypeUrl      string             `json:"type_url,omitempty",yaml:"type_url,omitempty"`
	Nonce        string             `json:"nonce,omitempty",yaml:"nonce,omitempty"`
	ControlPlane *core.ControlPlane `json:"control_plane,omitempty",yaml:"control_plane,omitempty"`
}

type marshaller interface {
	marshal(*apiv2.DiscoveryResponse) (string, error)
}

type jsonMarshaller struct {
	keepIndent bool
}

type defaultMarshaller struct{}
type yamlMarshaller struct{}

func convertToStructuredDiscoveryResponse(raw *apiv2.DiscoveryResponse) (*discoveryResponse, error) {
	resp := &discoveryResponse{
		VersionInfo:  raw.GetVersionInfo(),
		Resources:    make([]interface{}, len(raw.GetResources())),
		Canary:       raw.GetCanary(),
		TypeUrl:      raw.GetTypeUrl(),
		Nonce:        raw.GetNonce(),
		ControlPlane: raw.GetControlPlane(),
	}

	for i, item := range raw.GetResources() {
		switch item.GetTypeUrl() {
		case _typeURLMap["eds"]:
			target := &apiv2.ClusterLoadAssignment{}
			if err := proto.Unmarshal(item.GetValue(), target); err != nil {
				return nil, err
			}
			resp.Resources[i] = target
		default:
			return nil, _errUnknownTypeUrl
		}
	}

	return resp, nil
}

func newJSONMarshaller() marshaller {
	return &jsonMarshaller{keepIndent: true}
}

func (f *jsonMarshaller) marshal(raw *apiv2.DiscoveryResponse) (string, error) {
	resp, err := convertToStructuredDiscoveryResponse(raw)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(resp)
	return string(data), err
}

func newDefaultMarshaller() marshaller {
	return &defaultMarshaller{}
}

func (f *defaultMarshaller) marshal(raw *apiv2.DiscoveryResponse) (string, error) {
	resp, err := convertToStructuredDiscoveryResponse(raw)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(*resp), nil
}

func newYAMLMarshaller() marshaller {
	return &yamlMarshaller{}
}

func (f *yamlMarshaller) marshal(raw *apiv2.DiscoveryResponse) (string, error) {
	resp, err := convertToStructuredDiscoveryResponse(raw)
	if err != nil {
		return "", err
	}
	data, err := yaml.Marshal(resp)
	return string(data), err
}
